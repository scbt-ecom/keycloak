package keycloak

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL  string
	ClientID string
	Realm    string
	Scope    string
	cl       *http.Client
}

var keycloakClient Client

func NewClient(baseURL string, clientID string, realm string, scope string) {
	keycloakClient = Client{
		BaseURL:  baseURL,
		ClientID: clientID,
		Realm:    realm,
		Scope:    scope,
		cl: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Proxy: nil,
			},
		},
	}
}
