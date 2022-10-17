package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/controllers"
	"github.com/panhdjf/server_management_system/middleware"
)

type ServerRouteController struct {
	serverController controllers.SeverController
}

func NewRouteServerController(serverController controllers.SeverController) ServerRouteController {
	return ServerRouteController{serverController}
}

func (sc ServerRouteController) ServerRoute(rg *gin.RouterGroup) {

	router := rg.Group("servers")
	router.Use(middleware.DeserializeUser())
	router.POST("/", sc.serverController.CreateServer)
	router.GET("/:serverIpv4", sc.serverController.ViewServers)
	router.PUT("/:serverId", sc.serverController.UpdateServer)
	router.DELETE("/:serverId", sc.serverController.DeletePost)
}

// 	router.DELETE("/", c.servercontroller.Delete_all_servers)
// 	router.POST("/all", c.servercontroller.CreatemanyServer)
// 	router.GET("/all/port", c.servercontroller.Check_on_off)
// 	router.GET("/ipv4/:ipv4", c.servercontroller.Check)
