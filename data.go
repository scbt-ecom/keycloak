package keycloak

import "net/http"

func isHaveAccessToken(req *http.Request) (string, bool) {
	token, err := req.Cookie("access_token")
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
