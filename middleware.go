package keycloak

import (
	"context"
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

func (cl *Client) AuthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	code, have := isHaveQueryCode(r)
	if !have {
		slog.Info("redirect to authorization page",
			slogging.StringAttr("url", r.URL.String()),
		)

		http.Redirect(w, r, generateCodeURL(cl), http.StatusFound)
		return
	}

	slog.Info("starting token request")
	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType: "client",
		code:        code,
	}, cl)
	if err != nil {
		if os.IsTimeout(err) {
			slog.Info("keycloak token request timed out")
			sendError(w, http.StatusGatewayTimeout, errNetworkAccess)
			return
		}

		slog.Info("keycloak token request failed with error",
			slogging.ErrAttr(err))
		sendError(w, http.StatusInternalServerError, err)
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
	errNoRoles            = errors.New("user dont have one of roles")
)

func (cl *Client) NeedRoleDirectRedirect(requiredRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if requiredRoles == nil {
				next.ServeHTTP(w, r)
				return
			}

			accessToken, have := isHaveAccessToken(r)
			if !have {
				slog.Info("redirect to authorization page",
					slogging.StringAttr("url", r.URL.String()),
				)

				http.Redirect(w, r, generateCodeURL(cl), http.StatusFound)
				return
			}

			userRoles, err := introspectTokenRoles(accessToken, cl.ClientID)
			if err != nil {
				slog.Error("failed to get user roles",
					slogging.ErrAttr(err))
				sendError(w, http.StatusInternalServerError, err)
				return
			}

			username, err := extractUsername(accessToken)
			if err != nil {
				slog.Error("failed to get username",
					slogging.ErrAttr(err))
				sendError(w, http.StatusInternalServerError, err)
				return
			}

			if !isHaveRole(userRoles, requiredRoles) {
				slog.Error("user dont have one of roles")
				sendError(w, http.StatusForbidden, errNoRoles)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "username", username)))
		})
	}
}

func (cl *Client) NeedRole(requiredRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if requiredRoles == nil {
				next.ServeHTTP(w, r)
				return
			}

			accessToken, have := isHaveAccessToken(r)
			if !have {
				sendError(w, http.StatusForbidden, errors.New("user has no access token in cookie"))
				return
			}

			userRoles, err := introspectTokenRoles(accessToken, cl.ClientID)
			if err != nil {
				slog.Error("failed to get user roles",
					slogging.ErrAttr(err))
				sendError(w, http.StatusInternalServerError, err)
				return
			}

			username, err := extractUsername(accessToken)
			if err != nil {
				slog.Error("failed to get username",
					slogging.ErrAttr(err))
				sendError(w, http.StatusInternalServerError, err)
				return
			}

			if !isHaveRole(userRoles, requiredRoles) {
				sendError(w, http.StatusForbidden, errNoRoles)
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
//				w.Write(beautifyError(err))
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
