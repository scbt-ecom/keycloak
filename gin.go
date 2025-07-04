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
		if requiredRoles == nil {
			c.Next()
			return
		}

		accessToken, have := isHaveAccessToken(c.Request)
		if !have {
			slogging.L(c.Request.Context()).Warn("user have no access token in cookie",
				slogging.StringAttr("url", c.Request.URL.String()))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userRoles, err := introspectTokenRoles(accessToken, cl.ClientID)
		if err != nil {
			slogging.L(c.Request.Context()).Warn("failed to get user roles",
				slogging.ErrAttr(err),
			)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		username, err := extractUsername(accessToken)
		if err != nil {
			slogging.L(c.Request.Context()).Warn("failed to extract username",
				slogging.ErrAttr(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !isHaveRole(userRoles, requiredRoles) {
			slogging.L(c.Request.Context()).Warn("user does not have required role",
				slogging.AnyAttr("required_roles", requiredRoles),
				slogging.AnyAttr("user_roles", userRoles))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c, "username", username))
		c.Next()
	}
}
