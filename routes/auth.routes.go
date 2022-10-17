package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/controllers"
	"github.com/panhdjf/server_management_system/middleware"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (ac *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", ac.authController.SignUpUser)
	router.POST("/login", ac.authController.SignInUser)
	router.GET("/refresh", ac.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(), ac.authController.LogoutUser)
}
