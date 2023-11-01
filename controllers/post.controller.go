package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dangtran47/go_crud/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(db *gorm.DB) PostController {
	return PostController{DB: db}
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreatePost

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	newPost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		AuthorID:  currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := pc.DB.Create(&newPost)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Post already exists"})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"error": result.Error.Error()})
		return
	}

	postReponse := newPost.ToResponse()
	ctx.JSON(http.StatusCreated, gin.H{"data": postReponse})
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdatePost
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedPost models.Post
	result := pc.DB.First(&updatedPost, "id = ? AND author_id = ?", postId, currentUser.ID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	now := time.Now()
	updatedPost.Title = payload.Title
	updatedPost.Content = payload.Content
	updatedPost.UpdatedAt = now

	pc.DB.Save(updatedPost)

	postReponse := updatedPost.ToResponse()

	ctx.JSON(http.StatusOK, gin.H{"data": postReponse})
}

func (pc *PostController) GetPost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var post models.Post
	result := pc.DB.First(&post, "id = ? AND author_id = ?", postId, currentUser.ID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	postReponse := post.ToResponse()

	ctx.JSON(http.StatusOK, gin.H{"data": postReponse})
}

func (pc *PostController) GetAllPosts(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var posts []models.Post
	result := pc.DB.Limit(limit).Offset(offset).Where("author_id = ?", currentUser.ID).Find(&posts)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	postResponses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	ctx.JSON(http.StatusOK, gin.H{"data": postResponses})
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	postId := ctx.Param("postId")

	result := pc.DB.Delete(&models.Post{}, "id = ? AND author_id = ?", postId, currentUser.ID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
