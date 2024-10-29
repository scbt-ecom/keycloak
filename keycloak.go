package keycloak

type Config struct {
	needExternalAuthorization bool
	externalSettings          *externalAuthorizationSettings
}

type externalAuthorizationSettings struct {
	baseURL  string
	realm    string
	clientID string
}

const (
	defaultNeedExternalAuthorization = false
)

//type KeycloakOption func(*Config)

//func NewConfig(opts ...KeycloakOption) *Config {
//	cfg := &Config{
//		needExternalAuthorization: defaultNeedExternalAuthorization,
//		externalSettings:          nil,
//	}
//
//	for _, opt := range opts {
//		opt(cfg)
//	}
//
//	return cfg
//}
//
//func createGetTokenRequest(code string, redirectURL string) (*http.Request, error) {
//	tokenUrl := fmt.Sprintf("%sauth/realms/%s/protocol/openid-connect/token", keycloakClient.BaseURL, keycloakClient.Realm)
//
//	data := url.Values{}
//	data.Set("grant_type", "authorization_code")
//	data.Set("client_id", keycloakClient.ClientID)
//	data.Set("code", code)
//	data.Set("redirect_uri", redirectURL)
//
//	req, err := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(data.Encode()))
//	if err != nil {
//		return nil, err
//	}
//
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//	return req, nil
//}

//func DoTokenRequest(req *http.Request) (accessToken string, err error) {
//	resp, err := keycloakClient.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	bb, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//
//	var tokenResp TokenResponse
//	err = json.Unmarshal(bb, &tokenResp)
//	if err != nil {
//		return "", err
//	}
//
//	if tokenResp.AccessToken == "" {
//		return "", errors.New("no access token found")
//	}
//
//	return tokenResp.AccessToken, nil
//}
