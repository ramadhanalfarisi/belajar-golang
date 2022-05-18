package routers

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Router() {
	mux := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET","POST","PUT","DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	ed_version := mux.PathPrefix("/v1").Subrouter()
	// secure := ed_version.NewRoute().Subrouter()
	Authentication(ed_version)
	http.ListenAndServe(":3000",handlers.CORS(headers,methods,origins)(mux))
}