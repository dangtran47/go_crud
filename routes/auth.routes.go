package routes

import (
	"github.com/dangtran47/go_crud/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	route := rg.Group("/auth")

	route.POST("/signup", rc.authController.SignUpUser)
	route.POST("/signin", rc.authController.SignInUser)
	route.GET("/refresh", rc.authController.RefreshAccessToken)
	route.GET("/signout", rc.authController.SignOut)
	route.GET("/verify/:code", rc.authController.VerifyEmail)
}
