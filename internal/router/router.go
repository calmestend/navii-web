package router

import (
	"net/http"

	"github.com/naviiapp/web/internal/handlers"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("GET /user/delete", handlers.DeleteAccountForm)
	mux.HandleFunc("POST /user/delete", handlers.DeleteAccount)
	mux.HandleFunc("GET /privacy_policies", handlers.PrivacyPolicies)

	return mux
}
