package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var db_user, db_password, db_name, db_host, db_port string

func Migrate(environment string) {
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
	strCon := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, db_host, db_port, db_name)
	db, err := sql.Open("postgres", strCon)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	path := "file://./migrations/"
	m, err := migrate.NewWithDatabaseInstance(
		path,
		db_name, driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Println(err)
	}
}
