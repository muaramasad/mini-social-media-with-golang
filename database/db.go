package database

import (
	//"fmt"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *gorm.DB {

	db, err := gorm.Open("mysql", "go_user:go_password@tcp(127.0.0.1:4406)/go_twitter_db?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}

	return db
}