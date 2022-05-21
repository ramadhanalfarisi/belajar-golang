package routers

import (
	"belajar_golang/controllers"

	"github.com/gorilla/mux"
)

func Products(mux *mux.Router) {
	mux.HandleFunc("/products", controllers.SelectAllProducts).Methods("GET")
}