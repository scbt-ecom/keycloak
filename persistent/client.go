package persistent

import (
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseURL string
	realm   string
	*http.Client
}

func NewClient(baseURL string, realm string) (*Client, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		return nil, fmt.Errorf("base url can't be empty")
	}

	if realm == "" {
		return nil, fmt.Errorf("realm can't be empty")
	}

	return &Client{
		baseURL: baseURL,
		realm:   realm,
		Client:  http.DefaultClient,
	}, nil
}
