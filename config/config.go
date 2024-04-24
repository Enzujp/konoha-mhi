package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/caarlos0/env.v2"
)


type Configuration struct {
	Development	bool `env:"DEVELOPMENT" envDefault:"true"`
	Port	int	`env:"PORT" required:"true"`
	DB_URL	string	`env:"DB_URL" required:"true"`
}

func GetEnv() (*Configuration, error) {
	var config Configuration
	// Parse env values to access their values
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	if config.Development {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		portString := os.Getenv("PORT")
		if portString == "" {
			log.Fatal("Empty Port String")
	}
		log.Println(".env file loaded successfully")
		if err := env.Parse(&config); err != nil {
			log.Fatal(err.Error())
		}
	}
	return &config, nil
}
