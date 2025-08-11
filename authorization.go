package keycloak

import (
	"encoding/base64"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

type Credentials struct {
	ClientID     string
	ClientSecret string
}

type AuthData struct {
	accessToken string

	mu           sync.Mutex
	refreshToken string
}

func (a *AuthData) GetAccessToken() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.accessToken
}

func (cl *Client) AuthWithCredentials(creds Credentials) (*AuthData, error) {
	data, err := doTokenRequest(&tokenRequestData{
		requestType:  "",
		code:         "",
		clientID:     "",
		clientSecret: "",
	}, cl)
	if err != nil {
		return nil, err
	}

	accessExp, err := introspectTokenExp(data.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshExp, err := introspectTokenExp(data.RefreshToken)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	accessExpiresIn := time.Duration(accessExp-now) * time.Second
	refreshExpiresIn := time.Duration(refreshExp-now) * time.Second

	accessExpires := time.NewTicker(accessExpiresIn)
	defer accessExpires.Stop()

	refreshExpires := time.NewTicker(refreshExpiresIn)
	defer refreshExpires.Stop()

	for {
		select {
		case <-accessExpires.C:
		case <-refreshExpires.C:

		}
	}

}

const refreshTokenURL = "%s/realms/%s/protocol/openid-connect/token"

//func (cl *Client) RefreshToken(refreshToken string, creds Credentials) error {
//	url := fmt.Sprintf(refreshTokenURL, cl.BaseURL, cl.Realm)
//
//	formData := map[string]string{
//		"grant_type":    "refresh_token",
//		"client_id":     creds.ClientID,
//		"client_secret": creds.ClientSecret,
//		"refresh_token": refreshToken,
//	}
//
//	requestBody := bytes.NewBufferString("")
//	for k, v := range formData {
//		requestBody.WriteString(fmt.Sprintf("%s=%s&", k, v))
//	}
//
//	resp, err := http.Post(url, "application/x-www-form-urlencoded", requestBody)
//	if err != nil {
//		return fmt.Errorf("failed to send request: %s", err.Error())
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return fmt.Errorf("failed to read response body: %s", err.Error())
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed to refresh token: %s", string(body))
//	}
//
//	return nil
//}

func introspectTokenExp(token string) (int64, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return -1, errInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return -1, err
	}

	result := gjson.GetBytes(payload, "exp")
	return result.Int(), nil
}

//func (cl *Client) AuthWithCredentials(creds Credentials) (*AuthData, error) {
//
//	tokenData, err := doTokenRequest(&tokenRequestData{
//		requestType:  "server",
//		clientID:     creds.ClientID,
//		clientSecret: creds.ClientSecret,
//	}, cl)
//	if err != nil {
//		return nil, err
//	}
//
//	return &AuthData{
//		AccessToken:  tokenData.AccessToken,
//		RefreshToken: tokenData.RefreshToken,
//	}, nil
//}
