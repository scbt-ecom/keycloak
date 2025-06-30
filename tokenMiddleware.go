package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const userKey string = "user"

// This middleware is used with the token that comes in the Authorization header
func (cl *Client) NeedTokenRole(requiredRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := jwt.ParseHeader(r.Header, "Authorization", jwt.WithVerify(false), jwt.WithValidate(false))
			if err != nil {
				sendError(w, http.StatusUnauthorized, err)
				return
			}

			auth, err := NewAuthenticator(tokenStr, r.Context())
			if err != nil {
				sendError(w, http.StatusUnauthorized, err)
				return
			}

			if auth.Claims.Issuer == "" {
				sendError(w, http.StatusUnauthorized, errors.New("missing token issuer"))
				return
			}

			if len(auth.Claims.Audience) == 0 {
				sendError(w, http.StatusUnauthorized, errors.New("missing token audience"))
				return
			}

			now := time.Now()
			if auth.Claims.Expiration.Before(now) {
				sendError(w, http.StatusUnauthorized, errors.New("access token expired"))
				return
			}
			if auth.Claims.NotBefore.After(now) {
				sendError(w, http.StatusUnauthorized, errors.New("token not yet valid"))
				return
			}

			haveRole := isHaveRole(auth.Realm, requiredRoles)
			if !haveRole {
				haveRole = isHaveRole(joinResourceAccess(auth.Resource), requiredRoles)
				if !haveRole {
					sendError(w, http.StatusForbidden, errNoRoles)
					return
				}
			}

			userData := map[string]interface{}{
				"sub":             auth.Claims.Subject,
				"username":        auth.Claims.Username,
				"realm_roles":     auth.Realm,
				"resource_roles":  auth.Resource,
				"scope":           auth.Claims.Scope,
				"issuer":          auth.Claims.Issuer,
				"audience":        auth.Claims.Audience,
				"resource_access": auth.Claims.ResourceAccess,
			}
			ctx := context.WithValue(r.Context(), userKey, userData)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func sendError(w http.ResponseWriter, statusCode int, err error) {
	errMessage := map[string]string{
		"error": err.Error(),
	}
	data, _ := json.Marshal(errMessage)

	w.WriteHeader(statusCode)
	w.Write(data)
}
