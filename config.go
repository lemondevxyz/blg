package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/thanhpk/randstr"
)

type Config struct {
	DisableRegistration bool   // disable registration(requires restart)
	Domain              string // used for RSS
	Secret              string // session key, 32len
	Title               string
	Description         string
	Port                uint16
}

var config Config

func init() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("title", "John's blog")
	viper.SetDefault("description", "A blog about programming")
	viper.SetDefault("domain", "localhost:8080")
	viper.SetDefault("secret", randstr.String(32))
	viper.SetDefault("port", 8080)

	var err error

	// Initiate viper for our config
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			os.Create("config.yaml")
			viper.WriteConfig()
		}
	}
	log.Println(config)

	// Unmarshal the config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to unmarshal config, error: %v", err)
	}
	log.Println(config)

}
