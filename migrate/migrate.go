package main

import (
	"fmt"
	"log"

	"github.com/panhdjf/server_management_system/initializers"
	"github.com/panhdjf/server_management_system/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Server{})
	fmt.Println("Migration complete")
}
