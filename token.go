package keycloak

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

var (
	errInvalidToken          = errors.New("invalid token")
	errInvalidKeycloakConfig = errors.New("invalid keycloak config")
)

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
