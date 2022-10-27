package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/middleware"
)

func (sc *ServerRouteController) DailyReportManuallyRoute(rg *gin.RouterGroup) {

	router := rg.Group("dailyreport")
	router.Use(middleware.DeserializeUser())
	router.GET("/manually", sc.serverController.DailyReportManually)
}
