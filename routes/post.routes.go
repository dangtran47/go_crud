package routes

import (
	"github.com/dangtran47/go_crud/controllers"
	"github.com/dangtran47/go_crud/middleware"
	"github.com/gin-gonic/gin"
)

type PostRouteController struct {
	postRouteController controllers.PostController
}

func NewPostRouteController(postRouteController controllers.PostController) PostRouteController {
	return PostRouteController{postRouteController}
}

func (pc *PostRouteController) PostRoute(rg *gin.RouterGroup) {
	route := rg.Group("/posts")
	route.Use(middleware.DeserializeUser())
	route.POST("/", pc.postRouteController.CreatePost)
	route.PUT("/:postId", pc.postRouteController.UpdatePost)
	route.DELETE("/:postId", pc.postRouteController.DeletePost)
	route.GET("/:postId", pc.postRouteController.GetPost)
	route.GET("/", pc.postRouteController.GetAllPosts)
}
