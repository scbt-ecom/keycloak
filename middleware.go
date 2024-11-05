package keycloak

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/scbt-ecom/slogging"
	"log/slog"
	"net/http"
	"os"
)

var (
	errNetworkAccess = errors.New("problem with network access")
)

func AuthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	code, have := isHaveQueryCode(r)
	if !have {
		slog.Info("dont have query code",
			slogging.StringAttr("url", r.URL.String()))
		http.Redirect(w, r, generateCodeURL(keycloakClient.RedirectURL), http.StatusFound)
		return
	}

	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType: "client",
		code:        code,
	})
	if err != nil {
		slog.Info("error while token request",
			slogging.ErrAttr(err))

		if os.IsTimeout(err) {
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write(beatifyError(errNetworkAccess))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(beatifyError(err))
		return
	}

	setupCookie(w, tokenData)
	slog.Info("successfully setup cookie")

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	slog.Info("successfully redirect to /")
}

var (
	errSomethingWentWrong = errors.New("something went wrong")
	errStatusNotOK        = errors.New("external resource response status not OK")
	errInvalidRequest     = errors.New("invalid request type, contact with developer")
)

func NeedRole(requiredRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, have := isHaveAccessToken(r)
			if !have {
				http.Redirect(w, r, generateCodeURL(keycloakClient.RedirectURL), http.StatusMovedPermanently)
				return
			}

			userRoles, err := introspectTokenRoles(accessToken)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(beatifyError(err))
				return
			}

			if !isHaveRole(userRoles, requiredRoles) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

//func MuxNeedAllRoles(requiredRoles ...string) mux.MiddlewareFunc {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			accessToken, have := isHaveAccessToken(r)
//			if !have {
//				http.Redirect(w, r, generateCodeURL(keycloakClient.RedirectURL), http.StatusMovedPermanently)
//				return
//			}
//
//			userRoles, err := introspectTokenRoles(accessToken)
//			if err != nil {
//				w.WriteHeader(http.StatusInternalServerError)
//				w.Write(beatifyError(err))
//				return
//			}
//
//			if !isHaveAllRoles(userRoles, requiredRoles) {
//				w.WriteHeader(http.StatusForbidden)
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
//}

func beatifyError(err error) []byte {
	errMessage := map[string]string{
		"error": err.Error(),
	}

	data, _ := json.Marshal(errMessage)
	return data
}
