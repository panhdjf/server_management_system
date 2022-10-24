package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/panhdjf/server_management_system/controllers"
// 	"github.com/panhdjf/server_management_system/middleware"
// )

// type MailRouteController struct {
// 	mailController controllers.MailController
// }

// func NewRouteMailController(mailController controllers.MailController) MailRouteController {
// 	return MailRouteController{mailController}
// }

// func (mc *MailRouteController) MailRoute(rg *gin.RouterGroup) {

// 	router := rg.Group("mail")
// 	router.GET("/", middleware.DeserializeUser(), mc.mailController.Cron)
// }
