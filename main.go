package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func main() {

	validate = validator.New()

	var err error

	db, err = CreateDatabase("db.sqlite3")
	if err != nil {
		log.Fatalf("An error occured with initializing the database: %v", err)
	}

	s := &server{
		r:      gin.Default(),
		port:   config.Port,
		secret: config.Secret,
	}

	s.Start()
}
