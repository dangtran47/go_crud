package main

import (
	"log"
	"net/http"

	"github.com/dangtran47/go_crud/controllers"
	"github.com/dangtran47/go_crud/initializers"
	"github.com/dangtran47/go_crud/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthController := controllers.NewAuthController(initializers.DB)
	AuthRouteController := routes.NewAuthRouteController(AuthController)
	AuthRouteController.AuthRoute(router)

	UserController := controllers.NewUserController(initializers.DB)
	UserRouteController := routes.NewUserRouteController(UserController)
	UserRouteController.UserRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
