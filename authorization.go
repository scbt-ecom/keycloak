package keycloak

type Credentials struct {
	ClientID     string
	ClientSecret string
}

type AuthData struct {
	AccessToken  string
	RefreshToken string
}

func (cl *Client) AuthWithCredentials(creds Credentials) (*AuthData, error) {
	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType:  "server",
		clientID:     creds.ClientID,
		clientSecret: creds.ClientSecret,
	}, cl)
	if err != nil {
		return nil, err
	}

	return &AuthData{
		AccessToken:  tokenData.AccessToken,
		RefreshToken: tokenData.RefreshToken,
	}, nil
}
