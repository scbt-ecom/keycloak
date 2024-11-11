package keycloak

import (
	"fmt"
	"net/http"
	"time"
)

func isHaveAccessToken(req *http.Request) (string, bool) {
	token, err := req.Cookie("access_token")
	if err != nil {
		return "", false
	}

	return token.Value, true
}

func isHaveRefreshToken(req *http.Request) (string, bool) {
	token, err := req.Cookie("refresh_token")
	if err != nil {
		return "", false
	}

	return token.Value, true
}

func isHaveQueryCode(req *http.Request) (string, bool) {
	code := req.URL.Query().Get("code")
	if code == "" {
		return "", false
	}

	return code, true
}

func setupCookie(w http.ResponseWriter, tokenData *tokenResponseData) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenData.AccessToken,
		Path:     "/",
		HttpOnly: false,
		Expires:  time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenData.RefreshToken,
		Path:     "/",
		HttpOnly: false,
		Expires:  time.Now().Add(time.Duration(tokenData.RefreshExpiresIn) * time.Second),
	})
}

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

//func isHaveAllRoles(userRoles []string, requiredRoles []string) bool {
//	for _, requiredRole := range requiredRoles {
//		for _, userRole := range userRoles {
//			if userRole == requiredRole {
//				break
//			}
//		}
//		return false
//	}
//
//	return true
//}

func generateCodeURL(redirectURL string) string {
	return fmt.Sprintf("%sauth/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s", cl.BaseURL, cl.Realm, cl.ClientID, redirectURL, cl.Scope)
}
