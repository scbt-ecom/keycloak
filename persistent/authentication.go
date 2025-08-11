package persistent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/scbt-ecom/keycloak/v2"
)

const tokenURL = "%s/auth/realms/%s/protocol/openid-connect/token"

func (cl *Client) authWithCredentials(creds keycloak.Credentials) (*authData, error) {
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {creds.ClientID},
		"client_secret": {creds.ClientSecret},
	}

	preparedURL := fmt.Sprintf(tokenURL, cl.baseURL, cl.realm)

	req, err := http.NewRequest(http.MethodPost, preparedURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("bad response status code: %d", resp.StatusCode))
	}

	var respBody authData
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, fmt.Errorf("decode response body: %s", err.Error())
	}

	return &respBody, nil
}

//func introspectTokenExp(token string) int64 {
//	return gjson.Get(token, "expires_in").Int()
//}

type authData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
