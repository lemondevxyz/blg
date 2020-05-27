package main

import (
	"log"
	"os"
	"testing"

	"gopkg.in/go-playground/validator.v9"
)

func TestMain(m *testing.M) {
	var err error

	/*
		s := &server{secret: "secret", port: 8080, r: gin.New()}
		s.initAPI()
		r = s.r
	*/

	db, err = CreateDatabase("test.sqlite3")
	if err != nil {
		log.Fatalf("Error with creating the database: %v", err)
	}

	validate = validator.New()

	code := m.Run()
	os.Remove("test.sqlite3")
	os.Exit(code)
}
