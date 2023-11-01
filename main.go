package main

import (
	"log"
	"net/http"

	"github.com/dangtran47/go_crud/controllers"
	_ "github.com/dangtran47/go_crud/docs"
	"github.com/dangtran47/go_crud/initializers"
	"github.com/dangtran47/go_crud/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	server *gin.Engine

	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController
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
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	PostController := controllers.NewPostController(initializers.DB)
	PostRouteController := routes.NewPostRouteController(PostController)
	PostRouteController.PostRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
