package keycloak

import (
	"log/slog"
	"net/http"
	"net/url"
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
	proxy, err := url.Parse("https://trofimovaa2:TheSanekTrof123!@10.80.96.28:9090")
	if err != nil {
		slog.Error("error while proxy")
	}

	keycloakClient = Client{
		BaseURL:     baseURL,
		ClientID:    clientID,
		Realm:       realm,
		Scope:       scope,
		RedirectURL: redirectURL,
		cl: &http.Client{
			Timeout: time.Second * 180,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		},
	}
}
