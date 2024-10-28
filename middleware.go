package keycloak

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func NeedRoles(requiredRoles ...string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					code := r.URL.Query().Get("code")
					if code == "" {
						w.Header().Set("Referrer-Policy", "strict-origin")
						w.Header().Set("Content-Type", "text/html")
						http.Redirect(w, r, generateCodeURL(fmt.Sprintf("https://%s", r.Host)), http.StatusFound)
						return
					} else {
						req, err := createGetTokenRequest(code, r.URL.Query().Get("redirect_uri"))
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write(beatifyError(err))
							return
						}

						accessToken, err := DoTokenRequest(req)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write(beatifyError(err))
							return
						}

						http.SetCookie(w, &http.Cookie{
							Name:  "token",
							Value: accessToken,
							Path:  "/",
						})
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(beatifyError(err))
					return
				}
			}

			token, _ = r.Cookie("token")

			roles, err := introspectTokenRoles(token.Value)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(beatifyError(err))
				return
			}

			for _, requiredRole := range requiredRoles {
				for _, role := range roles {
					if role == requiredRole {
						next(w, r)
						return
					}
				}
			}

			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}

func beatifyError(err error) []byte {
	errMessage := map[string]string{
		"error": err.Error(),
	}

	data, _ := json.Marshal(errMessage)
	return data
}
