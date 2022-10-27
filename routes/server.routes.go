package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/controllers"
	"github.com/panhdjf/server_management_system/middleware"
)

type ServerRouteController struct {
	serverController controllers.ServerController
}

func NewRouteServerController(serverController controllers.ServerController) ServerRouteController {
	return ServerRouteController{serverController}
}

func (sc ServerRouteController) ServerRoute(rg *gin.RouterGroup) {

	router := rg.Group("servers")
	router.Use(middleware.DeserializeUser())
	router.POST("/", sc.serverController.CreateServer)
	router.GET("/view", sc.serverController.ViewServers)
	router.PUT("/:serverId", sc.serverController.UpdateServer)
	router.DELETE("/:serverId", sc.serverController.DeleteServer)
	router.DELETE("/", sc.serverController.DeleteAllServers)
	router.POST("/excel/import", sc.serverController.ImportExcel)
	router.GET("/excel/export", sc.serverController.ExportExcel)
}
