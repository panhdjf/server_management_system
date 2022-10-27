package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/controllers"
	"github.com/panhdjf/server_management_system/initializers"
	"github.com/panhdjf/server_management_system/routes"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	ServerController      controllers.ServerController
	ServerRouteController routes.ServerRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	ServerController = controllers.NewServerController(initializers.DB)
	ServerRouteController = routes.NewRouteServerController(ServerController)

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

	go func() {
		ServerController.DailyReport()
	}()

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	ServerRouteController.ServerRoute(router)
	ServerRouteController.DailyReportManuallyRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))

	// go func() {
	// 	for {
	// 		// Sleep
	// 		// Get server list
	// 		// Each server:
	// 			// Call API of agent
	// 				// if failed => server off
	// 				// if successed => server on
	// 			// Get up time from reponse
	// 			// Update to database
	// 	}
	// }

	// GET /status :8000
	// Get server uptime
	// Return
}
