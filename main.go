package main

import (
	"log"
	"tokokocak/app"
	"tokokocak/helpers"
)

func main() {
	
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Your application running on http://localhost:8080")
	helpers.Migrate("dev")
	app := app.App{}
	app.Connection("development")
	app.Run(":8080")
}
