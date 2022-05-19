package controllers

import (
	"belajar-golang/helpers"
	"log"
	"net/http"
)

func SelectAllProducts(w http.ResponseWriter, r *http.Request){
	db, err := helpers.Connection()
	if err != nil {
		log.Fatal(err)
	}
	
	
}