package routers

import (
	"tokokocak/controllers"

	"github.com/gorilla/mux"
)

func(routers *Routers) Authentication(mux *mux.Router){
	controller := controllers.Controller{}
	controller.DB = routers.DB
	mux.HandleFunc("/login", controller.Login).Methods("POST")
	mux.HandleFunc("/register", controller.Register).Methods("POST")
}