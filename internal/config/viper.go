package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()
	err := godotenv.Load("/Users/efishery/Documents/socmed_golang/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	err = config.ReadInConfig()

	if err != nil {
		log.Fatalf("Error read config: %v", err)
	}

	return config
}
