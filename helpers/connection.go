package helpers

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func Connection() (*gorm.DB, error) {
	var db_user, db_password, db_name, db_host, db_port string

	if environment := os.Getenv("ENVIRONMENT"); environment == "test" {
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

	cstring := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", db_host, db_user, db_password, db_name, db_port)
	conn, err := gorm.Open(postgres.Open(cstring), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return conn, nil
}
