package main

import (
	"io/ioutil"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func CreateDatabase(name string) (*gorm.DB, error) {

	db, err := gorm.Open("sqlite3", name)

	if err != nil {
		log.Fatalf("An error occured with initializing the database, error: %v", err)
	}

	db.SetLogger(log.New(ioutil.Discard, "\r\n", 0))

	db.AutoMigrate(&Post{})
	db.AutoMigrate(&User{})

	return db, err
}

func CloseDatabase(db *gorm.DB) {
	db.Close()
}
