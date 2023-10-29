package controllers

import (
	"net/http"

	"github.com/dangtran47/go_crud/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("current_user").(models.User)

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Role:      currentUser.Role,
		Photo:     currentUser.Photo,
		Provider:  currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": userResponse})
}
