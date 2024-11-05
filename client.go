package keycloak

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL     string
	ClientID    string
	Realm       string
	Scope       string
	RedirectURL string
	cl          *http.Client
}

var keycloakClient Client

func NewClient(baseURL string, clientID string, realm string, scope string, redirectURL string) {
	keycloakClient = Client{
		BaseURL:     baseURL,
		ClientID:    clientID,
		Realm:       realm,
		Scope:       scope,
		RedirectURL: redirectURL,
		cl: &http.Client{
			Timeout: time.Second * 180,
		},
	}
}
