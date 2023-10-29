package routes

import (
	"github.com/dangtran47/go_crud/controllers"
	"github.com/dangtran47/go_crud/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	route := rg.Group("/users")
	route.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
}
