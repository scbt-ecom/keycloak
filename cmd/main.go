package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/scbt-ecom/keycloak"
	"net/http"
)

func main() {
	keycloak.NewClient(
		"https://keycloak-int-test.sovcombank.group/",
		"web-ecom",
		"office",
		"openid",
	)

	needRoles := keycloak.MuxNeedRoles("ecom-k8s")

	r := mux.NewRouter()

	sub := r.PathPrefix("/test/").Subrouter()
	sub.Use(needRoles.Middleware)

	sub.HandleFunc("/sad", sad)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		fmt.Println(err)
	}
}

func sad(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
