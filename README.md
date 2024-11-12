# Getting #
Install with indicating the version, нексус кривой и не совсем работает корректно когда версия пакета с гитхаба выходит за >=2.x.x, а я этого не знал, а он уже закешировал хай версии кеш которых в нексусе я не могу инвалидировать, поэтому пожалуйста устанавливайте указывая версию, иначе могут быть проблемы, актуальная версия ниже
```
go get github.com/scbt-ecom/keycloak@v1.9.21
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
http.HandleFunc("/auth", keycloak.AuthHandlerFunc)

enricherRoles := keycloak.NeedRole("exampleRole1", "exampleRole2")
http.Handle("/rules", enricherRoles(http.HandlerFunc(sad)))

http.ListenAndServe(":8080", nil)
```
#### Mux ####
```
r := mux.NewRouter()
r.HandleFunc("/auth", keycloak.AuthHandlerFunc)

rules := r.Path("/rules").Subrouter()
rules.Handle("/", ruleGetExampleHandler)

rules.Use(keycloak.NeedRole("exampleRole1", "exampleRole2"))
```
#### Gin ####
```
r := gin.Default()
r.Handle(http.MethodGet, "/auth", keycloak.GinAuthHandlerFunc)

rules := r.Group("/rules")
rules.Handle(http.MethodGet, "/", ruleGetExampleHandler)

rules.Use(keycloak.GinNeedRole("exampleRole1", "exampleRole2"))
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
