package keycloak

import "strings"

func fixURL(url string) string {
	return strings.TrimRight(url, "/")
}
