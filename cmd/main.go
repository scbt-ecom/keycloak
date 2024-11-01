package main

//import (
//	"fmt"
//	"github.com/gorilla/mux"
//	"github.com/scbt-ecom/keycloak"
//	"net/http"
//)
//
//func main() {
//	keycloak.NewClient(
//		"https://keycloak-int-test.sovcombank.group/",
//		"web-ecom",
//		"office",
//		"openid",
//		"https://test-ecom-internal-enricher-k8s.sovcombank.group",
//	)
//
//	r := mux.NewRouter()
//
//	auth := keycloak.NeedRole("ecom-k8s")
//
//	sub := r.PathPrefix("/test/").Subrouter()
//	sub.Use(auth.Middleware)
//
//	sub.Handle("/sad", keycloak.NeedRole("sad")(http.HandlerFunc(sad)))
//
//	err := http.ListenAndServe(":8081", r)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	http.HandleFunc("adsdas", keycloak.AuthHandlerFunc)
//
//	roles := keycloak.NeedRole("dasdas")
//	http.HandleFunc("asdsafasfasfas", keycloak.AuthHandlerFunc)
//	http.Handle("adsdas", roles(http.HandlerFunc(sad)))
//
//}
//
//func sad(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("Hello world"))
//}
//
//func (s *Service) NoteHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
//	notes, err := noteService.GetAllNotes(ctx)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Fprintln(w, "Note:", notes)
//}
