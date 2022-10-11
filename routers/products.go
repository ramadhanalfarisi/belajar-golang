package routers

import (
	"tokokocak/controllers"

	"github.com/gorilla/mux"
)

func(routers *Routers) Products(mux *mux.Router) {
	controller := controllers.Controller{}
	controller.DB = routers.DB
	mux.HandleFunc("/products", controller.SelectAllProducts).Methods("GET")
	mux.HandleFunc("/products/{id}", controller.SelectOneProduct).Methods("GET")
	mux.HandleFunc("/products", controller.InsertProducts).Methods("POST")
	mux.HandleFunc("/products/{id}", controller.UpdateProduct).Methods("PUT")
	mux.HandleFunc("/products/{id}", controller.DeleteProduct).Methods("DELETE")

}