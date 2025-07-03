# Getting #
```
go env -w GONOPROXY=github.com/scbt-ecom/*
go get -u github.com/scbt-ecom/keycloak/v2
```

# Usage #
## Initialize ##

### Client ###
This is example values. The keycloak should be initialized where it will be used or passed as a function parameter. 
```
keycloakClient := keycloak.NewClient(keycloak.Config{
		BaseURL:     "https://keycloak-int-test.sovcombank.group/",
		ClientID:    "web-ecom",
		Realm:       "office",
		RedirectURL: "https://test-ecom-internal-enricher-k8s.sovcombank.group/auth",
	})
```

### Server ###
```
keycloakServer := keycloak.NewClient(keycloak.Config{
        BaseURL:  cfg.KeycloakBaseURL,
        ClientID: cfg.KeycloakInternalAuthUsername,
        Realm:    cfg.KeycloakRealm,
    })
```
## Authorization ##
### Client ###
#### Native ####
```
http.HandleFunc("/auth", keycloakClient.AuthHandlerFunc)

enricherRoles := keycloakClient.NeedRole("exampleRole1", "exampleRole2")
http.Handle("/rules", enricherRoles(http.HandlerFunc(sad)))

http.ListenAndServe(":8080", nil)
```
#### Mux ####
```
r := mux.NewRouter()
r.HandleFunc("/auth", keycloakClient.AuthHandlerFunc)

rules := r.Path("/rules").Subrouter()
rules.Handle("/", ruleGetExampleHandler)

rules.Use(keycloakClient.NeedRole("exampleRole1", "exampleRole2"))
```
#### Gin ####
```
r := gin.Default()
r.Handle(http.MethodGet, "/auth", keycloakClient.GinAuthHandlerFunc)

rules := r.Group("/rules")
rules.Handle(http.MethodGet, "/", ruleGetExampleHandler)

rules.Use(keycloakClient.GinNeedRole("exampleRole1", "exampleRole2"))
```
### Server ###
```
authData, err := keycloakClient.AuthWithCredentials(
    keycloak.Credentials{
	    ClientID     : {exampleClientID},
	    ClientSecret : {exampleClientSecret},
    },
)
if err != nil {
    return nil, err
}
```

### TokenClient ###
It should be used if you only need token authentication from another server (not global keycloak from environment variables).
For example, for systems like api2api.  
```
r := mux.NewRouter()
r.HandleFunc("/auth", keycloakClient.AuthHandlerFunc)

tc := &keycloak.TokenClient{}

api := r.PathPrefix("/api").Subrouter()
api.Use(tc.NeedTokenRole("exampleRole1", "exampleRole1"))
api.HandleFunc("/", getExampleHandler)
```
