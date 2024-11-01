package keycloak

type Credentials struct {
	ClientID     string
	ClientSecret string
}

type AuthData struct {
	AccessToken  string
	RefreshToken string
}

func AuthWithCredentials(creds Credentials) (*AuthData, error) {
	tokenData, err := doTokenRequest(&tokenRequestData{
		requestType:  "server",
		clientID:     creds.ClientID,
		clientSecret: creds.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	return &AuthData{
		AccessToken:  tokenData.AccessToken,
		RefreshToken: tokenData.RefreshToken,
	}, nil
}
