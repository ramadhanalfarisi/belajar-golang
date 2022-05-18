package helpers

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connection() (*gorm.DB, error) {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	cstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db_user, db_password, db_host, db_port, db_name)
	conn, err := gorm.Open(mysql.Open(cstring), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return conn, nil
}