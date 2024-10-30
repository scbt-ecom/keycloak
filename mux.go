package keycloak

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultScheme = "http"

func dasdas(requiredRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := isHaveAccessToken(r)
			if !ok {
				//redirectURL, err := ensureRedirectURL(r)
				//if err != nil {
				//	w.WriteHeader(http.StatusInternalServerError)
				//	w.Write(beatifyError(err))
				//	return
				//}

				code, ok := isHaveQueryCode(r)
				if !ok {
					http.Redirect(w, r, generateCodeURL(keycloakClient.RedirectURL), http.StatusFound)
					return
				}

				accessToken, err := doTokenRequest(code)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(beatifyError(err))
					return
				}

				token = accessToken
			}

			userRoles, err := introspectTokenRoles(token)
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

func generateCodeURL(redirectURL string) string {
	return fmt.Sprintf("%sauth/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s", keycloakClient.BaseURL, keycloakClient.Realm, keycloakClient.ClientID, redirectURL, keycloakClient.Scope)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	errSomethingWentWrong = errors.New("something went wrong")
	errStatusNotOK        = errors.New("external resource response status not OK")
)

func doTokenRequest(code string) (string, error) {
	tokenURL := fmt.Sprintf("%sauth/realms/%s/protocol/openid-connect/token", keycloakClient.BaseURL, keycloakClient.Realm)

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"client_id":    {keycloakClient.ClientID},
		"code":         {code},
		"redirect_uri": {keycloakClient.RedirectURL},
	}

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := keycloakClient.cl.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errStatusNotOK
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(bb, &tokenResponse)
	if err != nil {
		return "", err
	}

	if tokenResponse.AccessToken == "" {
		return "", errSomethingWentWrong
	}

	return tokenResponse.AccessToken, nil
}

//func ensureRedirectURL(req *http.Request) (string, error) {
//	fullURL := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)
//
//	return fullURL, nil
//}

func isHaveRole(userRoles []string, requiredRoles []string) bool {
	for _, userRole := range userRoles {
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}

	return false
}

func MuxNeedRoles(requiredRoles ...string) mux.MiddlewareFunc {
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

func MuxAuthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	code, have := isHaveQueryCode(r)
	if !have {
		http.Redirect(w, r, generateCodeURL(keycloakClient.RedirectURL), http.StatusFound)
		return
	}

	accessToken, err := doTokenRequest(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(beatifyError(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
