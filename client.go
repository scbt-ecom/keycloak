package keycloak

import (
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL  string
	ClientID string
	Realm    string
	Scope    string
	*http.Client
}

var keycloakClient Client

func NewClient(baseURL string, clientID string, realm string, scope string) {
	keycloakClient = Client{BaseURL: baseURL, ClientID: clientID, Realm: realm, Scope: scope, Client: http.DefaultClient}
}

func generateCodeURL(redirectURL string) string {
	return fmt.Sprintf("%sauth/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s", keycloakClient.BaseURL, keycloakClient.Realm, keycloakClient.ClientID, encodedRedirectURL, keycloakClient.Scope)
}
