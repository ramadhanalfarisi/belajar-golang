package routers

import (
	"tokokocak/controllers"

	"github.com/gorilla/mux"
)

func Authentication(mux *mux.Router){
	mux.HandleFunc("/login", controllers.Login).Methods("POST")
	mux.HandleFunc("/register", controllers.Register).Methods("POST")
}