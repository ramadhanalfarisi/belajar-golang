package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tokokocak/middlewares"
	"tokokocak/routers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Mux *mux.Router
	DB  *gorm.DB
}

func (app *App) ListRouter() {
	mux := mux.NewRouter().StrictSlash(true)
	ed_version := mux.PathPrefix("/v1").Subrouter()
	secure := ed_version.NewRoute().Subrouter()
	secure.Use(middlewares.AuthMiddleware, middlewares.ApiMiddleware)
	app.Mux = mux
	routers := routers.Routers{}
	routers.DB = app.DB
	routers.Authentication(ed_version)
	routers.Products(secure)
}

func (app *App) Run(port string) {
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(port, handlers.CORS(headers, methods, origins)(app.Mux))
}

func (app *App) Connection(environment string) {
	var db_user, db_password, db_name, db_host, db_port string

	if environment == "test" {
		db_user = os.Getenv("DB_USER_TEST")
		db_password = os.Getenv("DB_PASSWORD_TEST")
		db_name = os.Getenv("DB_NAME_TEST")
		db_host = os.Getenv("DB_HOST_TEST")
		db_port = os.Getenv("DB_PORT_TEST")
	} else if environment == "development" {
		db_user = os.Getenv("DB_USER_DEV")
		db_password = os.Getenv("DB_PASSWORD_DEV")
		db_name = os.Getenv("DB_NAME_DEV")
		db_host = os.Getenv("DB_HOST_DEV")
		db_port = os.Getenv("DB_PORT_DEV")
	} else {
		db_user = os.Getenv("DB_USER")
		db_password = os.Getenv("DB_PASSWORD")
		db_name = os.Getenv("DB_NAME")
		db_host = os.Getenv("DB_HOST")
		db_port = os.Getenv("DB_PORT")
	}

	cstring := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", db_host, db_user, db_password, db_name, db_port)
	conn, err := gorm.Open(postgres.Open(cstring), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	app.DB = conn
	app.ListRouter()
}
