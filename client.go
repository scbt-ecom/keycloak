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

func NewClient(
	baseURL string,
	clientID string,
	realm string,
	scope string,
	redirectURL string,
) {
	keycloakClient = Client{
		BaseURL:     baseURL,
		ClientID:    clientID,
		Realm:       realm,
		Scope:       scope,
		RedirectURL: redirectURL,
		cl: &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				Proxy:                  nil,
				OnProxyConnectResponse: nil,
				DialContext:            nil,
				Dial:                   nil,
				DialTLSContext:         nil,
				DialTLS:                nil,
				TLSClientConfig:        nil,
				TLSHandshakeTimeout:    0,
				DisableKeepAlives:      false,
				DisableCompression:     false,
				MaxIdleConns:           0,
				MaxIdleConnsPerHost:    0,
				MaxConnsPerHost:        0,
				IdleConnTimeout:        0,
				ResponseHeaderTimeout:  0,
				ExpectContinueTimeout:  0,
				TLSNextProto:           nil,
				ProxyConnectHeader:     nil,
				GetProxyConnectHeader:  nil,
				MaxResponseHeaderBytes: 0,
				WriteBufferSize:        0,
				ReadBufferSize:         0,
				ForceAttemptHTTP2:      false,
			},
		},
	}
}
