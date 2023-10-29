package main

import (
	"fmt"
	"log"

	"github.com/dangtran47/go_crud/initializers"
	"github.com/dangtran47/go_crud/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	fmt.Println("AutoMigrate successfully")
}
