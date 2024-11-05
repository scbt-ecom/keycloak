package keycloak

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func sad(w *http.ResponseWriter, r *http.Request) {}

func GinAuthHandlerFunc(c *gin.Context) {
	code, have := isHaveQueryCode(c.Request)
	if !have {
		c.Redirect(http.StatusFound, generateCodeURL(cl.RedirectURL))
		return
	}

	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType: "client",
		code:        code,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	setupCookie(c.Writer, tokenData)

	c.Redirect(http.StatusMovedPermanently, "/")
}

func GinNeedRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, have := isHaveAccessToken(c.Request)
		if !have {
			c.Redirect(http.StatusMovedPermanently, generateCodeURL(cl.RedirectURL))
			return
		}

		userRoles, err := introspectTokenRoles(accessToken)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !isHaveRole(userRoles, requiredRoles) {
			c.Status(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
