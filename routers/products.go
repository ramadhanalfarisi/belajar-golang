package routers

import (
	"tokokocak/controllers"

	"github.com/gorilla/mux"
)

func(routers *Routers) Products(mux *mux.Router) {
	controller := controllers.Controller{}
	controller.DB = routers.DB
	mux.HandleFunc("/products", controller.SelectAllProducts).Methods("GET")
	mux.HandleFunc("/products/{id}", controller.SelectAllProducts).Methods("GET")
}