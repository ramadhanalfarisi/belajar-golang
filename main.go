package main

import (
	"tokokocak/routers"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "toko-kocak"
func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
    currentWorkDirectory, _ := os.Getwd()
    rootPath := projectName.Find([]byte(currentWorkDirectory))

    err := godotenv.Load(string(rootPath) + `/.env`)

    if err != nil {
        log.Fatalf("Error loading .env file")
    }
}

func main(){
	loadEnv()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Your application running on http://localhost:3000")
	routes := routers.Routes{}
	routes.Router()
}
