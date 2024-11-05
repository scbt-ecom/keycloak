package keycloak

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/scbt-ecom/slogging"
	"github.com/tidwall/gjson"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

var (
	errInvalidToken          = errors.New("invalid token")
	errInvalidKeycloakConfig = errors.New("invalid keycloak config")
)

type tokenResponseData struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  string `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type tokenRequestData struct {
	// client/server
	requestType string

	// code for client request
	code string

	// credentials for server request
	clientID     string
	clientSecret string
}

func doTokenRequest(reqData *tokenRequestData) (*tokenResponseData, error) {
	data := url.Values{}

	switch reqData.requestType {
	case "client":
		data = url.Values{
			"grant_type":   {"authorization_code"},
			"client_id":    {keycloakClient.ClientID},
			"code":         {reqData.code},
			"redirect_uri": {keycloakClient.RedirectURL},
		}
	case "server":
		data = url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {reqData.clientID},
			"client_secret": {reqData.clientSecret},
		}
	default:
		return nil, errInvalidRequest
	}

	tokenURL := fmt.Sprintf("%sauth/realms/%s/protocol/openid-connect/token", keycloakClient.BaseURL, keycloakClient.Realm)

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	slog.Info("keycloak token request started now")

	slog.Info("keycloak token request",
		slogging.AnyAttr("request1", req),
		slogging.AnyAttr("request2", &req),
		slogging.AnyAttr("request3", fmt.Sprintf("%+v", req)),
		slogging.AnyAttr("request4", fmt.Sprintf("%v", req)),
		slogging.AnyAttr("request5", fmt.Sprintf("%+v", &req)),
		slogging.AnyAttr("request6", fmt.Sprintf("%v", &req)),
	)

	slog.Info("keycloak token request",
		slogging.StringAttr("request11", fmt.Sprintf("%+v", *req)),
		slogging.StringAttr("request12", fmt.Sprintf("%v", *req)),

		slogging.StringAttr("request7", fmt.Sprintf("%+v", req)),
		slogging.StringAttr("request8", fmt.Sprintf("%v", req)),
	)

	resp, err := keycloakClient.cl.Do(req)
	if err != nil {
		slog.Error("keycloak token request failed",
			slogging.ErrAttr(err))
		return nil, err
	}
	defer resp.Body.Close()
	slog.Info("keycloak token response",
		slogging.AnyAttr("response", &resp))

	if resp.StatusCode != http.StatusOK {
		return nil, errStatusNotOK
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("keycloak token response reading failed",
			slogging.StringAttr("response", string(bb)))
		return nil, err
	}

	var tokenData tokenResponseData
	err = json.Unmarshal(bb, &tokenData)
	if err != nil {
		return nil, err
	}

	if tokenData.AccessToken == "" {
		return nil, errSomethingWentWrong
	}

	return &tokenData, nil
}

func introspectTokenRoles(token string) ([]string, error) {
	payload, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	result := gjson.GetBytes(payload, fmt.Sprintf("resource_access.%s.roles", keycloakClient.ClientID))

	roles, ok := result.Value().([]string)
	if !ok {
		return nil, errInvalidKeycloakConfig
	}

	return roles, nil
}
