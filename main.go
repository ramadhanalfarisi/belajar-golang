package main

import (
	"log"
	"tokokocak/app"
	"tokokocak/helpers"
)

func main() {
	env := "development"
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Your application running on http://localhost:8082")
	helpers.Migrate(env)
	app := app.App{}
	app.Connection(env)
	app.Run(":8080")
}
