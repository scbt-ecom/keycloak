package keycloak

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL     string
	ClientID    string
	Realm       string
	RedirectURL string
	cl          *http.Client
}

type Config struct {
	BaseURL     string
	ClientID    string
	Realm       string
	RedirectURL string
}

var cl Client

func NewClient(
	cfg Config,
) {
	cl = Client{
		BaseURL:     fixURL(cfg.BaseURL),
		ClientID:    cfg.ClientID,
		Realm:       cfg.Realm,
		RedirectURL: fixURL(cfg.RedirectURL),
		cl: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}
