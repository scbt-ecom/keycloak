# Getting #
```
go get github.com/scbt-ecom/keycloak@v1.2.14
```

# Usage #
## Initialize ##

### Client ###
This is example values
```
keycloak.NewClient(keycloak.Config{
		BaseURL:     "https://keycloak-int-test.sovcombank.group/",
		ClientID:    "web-ecom",
		Realm:       "office",
		Scope:       "openid",
		RedirectURL: "https://test-ecom-internal-enricher-k8s.sovcombank.group/auth",
	})
```

### Server ###
```
keycloak.NewClient(keycloak.Config{
        BaseURL:  cfg.KeycloakBaseURL,
        ClientID: cfg.KeycloakInternalAuthUsername,
        Realm:    cfg.KeycloakRealm,
    })
```
## Authorization ##
### Client ###
#### Native ####
```
http.HandleFunc("/auth", AuthHandlerFunc)

enricherRoles := NeedRole("exampleRole1", "exampleRole2")
http.Handle("/rules", enricherRoles(http.HandlerFunc(sad)))

http.ListenAndServe(":8080", nil)
```
#### Mux ####
```
r := mux.NewRouter()
r.HandleFunc("/auth", AuthHandlerFunc)

rules := r.Path("/rules").Subrouter()
rules.Handle("/", ruleGetExampleHandler)

rules.Use(NeedRole("exampleRole1", "exampleRole2"))
```
#### Gin ####
```
r := gin.Default()
r.Handle(http.MethodGet, "/auth", GinAuthHandlerFunc)

rules := r.Group("/rules")
rules.Handle(http.MethodGet, "/", ruleGetExampleHandler)

rules.Use(GinNeedRole("exampleRole1", "exampleRole2"))
```
### Server ###
```
authData, err := keycloak.AuthWithCredentials(
    keycloak.Credentials{
	    ClientID     : {exampleClientID},
	    ClientSecret : {exampleClientSecret},
    },
)
if err != nil {
    return nil, err
}
```
