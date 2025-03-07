package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/scbt-ecom/slogging"
)

var (
	errNetworkAccess = errors.New("problem with network access")
)

func AuthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	code, have := isHaveQueryCode(r)
	if !have {
		slog.Info("redirect to authorization page",
			slogging.StringAttr("url", r.URL.String()),
		)

		http.Redirect(w, r, generateCodeURL(cl.RedirectURL), http.StatusFound)
		return
	}

	slog.Info("starting token request")
	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType: "client",
		code:        code,
	})
	if err != nil {
		if os.IsTimeout(err) {
			slog.Info("keycloak token request timed out")
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write(beatifyError(errNetworkAccess))
			return
		}

		slog.Info("keycloak token request failed with error",
			slogging.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(beatifyError(err))
		return
	}
	slog.Info("token request succeeded")

	setupCookie(w, tokenData)
	slog.Info("successfully setup cookie from keycloak response")

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	slog.Info("successfully redirect to /")
}

var (
	errSomethingWentWrong = errors.New("something went wrong")
	errStatusNotOK        = errors.New("external resource response status not OK")
	errInvalidRequest     = errors.New("invalid request type, contact with developer")
)

func NeedRoleDirectRedirect(requiredRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, have := isHaveAccessToken(r)
			if !have {
				slog.Info("redirect to authorization page",
					slogging.StringAttr("url", r.URL.String()),
				)

				http.Redirect(w, r, generateCodeURL(cl.RedirectURL), http.StatusFound)
				return
			}

			userRoles, err := introspectTokenRoles(accessToken)
			if err != nil {
				slog.Error("failed to get user roles",
					slogging.ErrAttr(err))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(beatifyError(err))
				return
			}

			username, err := extractUsername(accessToken)
			if err != nil {
				slog.Error("failed to get username",
					slogging.ErrAttr(err))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(beatifyError(err))
				return
			}

			if !isHaveRole(userRoles, requiredRoles) {
				slog.Error("user dont have one of roles")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "username", username)))
		})
	}
}

func NeedRole(requiredRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, have := isHaveAccessToken(r)
			if !have {
				slog.Info("user have no access token in cookie")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			userRoles, err := introspectTokenRoles(accessToken)
			if err != nil {
				slog.Error("failed to get user roles",
					slogging.ErrAttr(err))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(beatifyError(err))
				return
			}

			username, err := extractUsername(accessToken)
			if err != nil {
				slog.Error("failed to get username",
					slogging.ErrAttr(err))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(beatifyError(err))
				return
			}

			if !isHaveRole(userRoles, requiredRoles) {
				slog.Error("user dont have one of roles")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "username", username)))
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
