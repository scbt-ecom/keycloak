package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type ResourceAccess map[string][]string
type RealmAccess []string

type Authenticator struct {
	Claims   TokenClaims
	Realm    RealmAccess
	Resource ResourceAccess
	Scope    []string
}

type TokenClaims struct {
	Subject        string                 `json:"sub"`
	Issuer         string                 `json:"iss"`
	Audience       []string               `json:"aud"`
	Expiration     time.Time              `json:"exp"`
	NotBefore      time.Time              `json:"nbf"`
	Scope          string                 `json:"scope"`
	Username       string                 `json:"preferred_username"`
	ResourceAccess map[string]interface{} `json:"resource_access"`
	ClientId       string                 `json:"client_id"`
}

func NewAuthenticator(token jwt.Token, ctx context.Context) (*Authenticator, error) {
	claimsMap, err := token.AsMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %w", err)
	}

	claimsJSON, err := json.Marshal(claimsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	var claims TokenClaims
	if err = json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	claims.Subject = token.Subject()
	aud := token.Audience()
	if len(aud) > 0 {
		claims.Audience = aud
	}

	privateClaims := token.PrivateClaims()
	if scopeStr, ok := privateClaims["scope"].(string); ok {
		claims.Scope = scopeStr
	}
	if resourceAccess, ok := privateClaims["resource_access"].(map[string]interface{}); ok {
		claims.ResourceAccess = resourceAccess
	}
	if clientId, ok := privateClaims["client_id"].(string); !ok {
		claims.ClientId = clientId
	}

	jwksURL := fmt.Sprintf("%s/protocol/openid-connect/certs", claims.Issuer)
	cache := jwk.NewCache(context.Background())
	cache.Register(jwksURL, jwk.WithMinRefreshInterval(24*time.Hour))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	keyset, err := cache.Get(ctx, jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	bytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	_, err = jwt.Parse(bytes, jwt.WithKeySet(keyset), jwt.WithVerify(false))
	if err != nil {
		return nil, fmt.Errorf("invalid token signature: %w", err)
	}

	realmRoles := extractRealmRoles(claimsMap)
	resourceRoles := extractRoles(claimsMap)

	return &Authenticator{
		Claims:   claims,
		Realm:    realmRoles,
		Resource: resourceRoles,
		Scope:    strings.Split(claims.Scope, " "),
	}, nil
}

func extractRealmRoles(claims interface{}) []string {
	var realmRoles []string
	if claimsMap, ok := claims.(map[string]interface{}); ok {
		if realmAccess, ok := claimsMap["realm_access"].(map[string]interface{}); ok {
			if r, ok := realmAccess["roles"].([]interface{}); ok {
				for _, role := range r {
					if roleStr, ok := role.(string); ok {
						realmRoles = append(realmRoles, roleStr)
					}
				}
			}
		}
	}
	return realmRoles
}

func extractRoles(claims interface{}) map[string][]string {
	resourceRoles := make(map[string][]string)

	claimsMap, _ := claims.(map[string]interface{})
	resourceAccess, _ := claimsMap["resource_access"].(map[string]interface{})

	for clientID, clientAccess := range resourceAccess {
		clientMap, _ := clientAccess.(map[string]interface{})
		clientRoles, _ := clientMap["roles"].([]interface{})

		roles := make([]string, 0, len(clientRoles))
		for _, role := range clientRoles {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}

		if len(roles) > 0 {
			resourceRoles[clientID] = roles
		}
	}

	return resourceRoles
}

func joinResourceAccess(resourceAccess map[string][]string) []string {
	var roles []string
	for _, resourceRoles := range resourceAccess {
		roles = append(roles, resourceRoles...)
	}
	return roles
}
