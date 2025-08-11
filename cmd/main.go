package main

import (
	"fmt"

	"github.com/scbt-ecom/keycloak/v2"
	"github.com/scbt-ecom/keycloak/v2/persistent"
)

type Service struct {
	session *persistent.Session
}

func main() {
	cl, err := persistent.NewClient("https://keycloak-int-test.sovcombank.group", "internalApi")
	if err != nil {
		return
	}

	session := cl.NewSession(keycloak.Credentials{
		ClientID:     "backendecomReserveAccount",
		ClientSecret: "examplesecret",
	})

	svc := Service{
		session: session,
	}

	accessToken, err := svc.session.GetAccessToken()
	if err != nil {
		return
	}

	fmt.Println(accessToken)
}
