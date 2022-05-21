package routers

import (
	"tokokocak/controllers"

	"github.com/gorilla/mux"
)

func Products(mux *mux.Router) {
	mux.HandleFunc("/products", controllers.SelectAllProducts).Methods("GET")
	mux.HandleFunc("/products/{id}", controllers.SelectAllProducts).Methods("GET")
}