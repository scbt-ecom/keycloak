package keycloak

import (
	"github.com/gin-gonic/gin"
	"github.com/scbt-ecom/slogging"
	"log/slog"
	"net/http"
	"os"
)

func GinAuthHandlerFunc(c *gin.Context) {
	code, have := isHaveQueryCode(c.Request)
	if !have {
		slog.Info("redirect to authorization page",
			slogging.StringAttr("url", c.Request.URL.String()),
		)

		c.Redirect(http.StatusFound, generateCodeURL(cl.RedirectURL))
		return
	}

	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType: "client",
		code:        code,
	})
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

func GinNeedRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, have := isHaveAccessToken(c.Request)
		if !have {
			slog.Info("user have no access token in cookie")
			c.Redirect(http.StatusMovedPermanently, generateCodeURL(cl.RedirectURL))
			return
		}

		userRoles, err := introspectTokenRoles(accessToken)
		if err != nil {
			slog.Error("failed to get user roles",
				slogging.ErrAttr(err))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !isHaveRole(userRoles, requiredRoles) {
			slog.Error("user dont have one of roles",
				slogging.ErrAttr(err))
			c.Status(http.StatusForbidden)
			return
		}

		slog.Info("user has role, authorization successful")

		c.Next()
	}
}
