package keycloak

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/scbt-ecom/slogging"
)

func (cl *Client) GinAuthHandlerFunc(c *gin.Context) {
	code, have := isHaveQueryCode(c.Request)
	if !have {
		slog.Info("redirect to authorization page",
			slogging.StringAttr("url", c.Request.URL.String()),
		)

		c.Redirect(http.StatusFound, generateCodeURL(cl))
		return
	}

	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType: "client",
		code:        code,
	}, cl)
	if err != nil {
		if os.IsTimeout(err) {
			slog.Info("keycloak token request timed out")
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": errNetworkAccess,
			})
			return
		}

		slog.Info("keycloak token request failed with error",
			slogging.ErrAttr(err))

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	slog.Info("token request succeeded")

	setupCookie(c.Writer, tokenData)
	slog.Info("successfully setup cookie from keycloak response")

	c.Redirect(http.StatusMovedPermanently, "/")
	slog.Info("successfully redirect to /")
}

func (cl *Client) GinNeedRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, have := isHaveAccessToken(c.Request)
		if !have {
			slog.Info("user have no access token in cookie")
			c.Status(http.StatusForbidden)
			return
		}

		userRoles, err := introspectTokenRoles(accessToken, cl.ClientID)
		if err != nil {
			slog.Error("failed to get user roles",
				slogging.ErrAttr(err))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		username, err := extractUsername(accessToken)
		if err != nil {
			slog.Error("failed to get username",
				slogging.ErrAttr(err))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !isHaveRole(userRoles, requiredRoles) {
			slog.Error("user dont have one of roles")
			c.Status(http.StatusForbidden)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c, "username", username))
		c.Next()
	}
}
