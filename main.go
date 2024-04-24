package main

import (
	"log"

	"github.com/enzujp/konoha-mhi/api"
	"github.com/enzujp/konoha-mhi/config"
	"github.com/enzujp/konoha-mhi/database"
)


func main(){
	// Initialize database
	database.ConnectDB()
	defer database.CloseDatabase()

	// Initialize configuration
	config, err := config.GetEnv()
	if err != nil {
		log.Fatalf("Error encountered in setting up config: %v", err)
	}

	a := &api.API{
		Config: config,
	}

	log.Fatal(a.Serve())
}