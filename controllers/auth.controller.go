package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/dangtran47/go_crud/initializers"
	"github.com/dangtran47/go_crud/models"
	"github.com/dangtran47/go_crud/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewAuthController(db *gorm.DB) AuthController {
	return AuthController{DB: db}
}

// SignUpUser godoc
//
//	@Summary		Sign up
//	@Description	Sign up a new user given email, name, password, password confirmation, and photo
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.SignUp	true	"Sign up payload"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/auth/signup [post]
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUp

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirmation {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password and password confirmation does not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Email:     payload.Email,
		Name:      payload.Name,
		Password:  hashedPassword,
		Role:      "user",
		Verified:  false,
		Photo:     payload.Photo,
		Provider:  "email",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	utils.SendVerificationEmail(&newUser)

	message := "An email has been sent to " + newUser.Email + " with verification code."
	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignIn
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email or password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	if !user.Verified {
		utils.SendVerificationEmail(&user)

		ctx.JSON(http.StatusForbidden, gin.H{"error": "please verify your email"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	accessToken, err := utils.GenerateToken(user.ID, config.AccessTokenPrivateKey, config.AccessTokenExpiresIn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateToken(user.ID, config.RefreshTokenPrivateKey, config.RefreshTokenExpiresIn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	errorMessage := "Could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	var user models.User
	result := ac.DB.First(&user, sub)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	accessToken, err := utils.GenerateToken(user.ID, config.AccessTokenPrivateKey, config.AccessTokenExpiresIn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (ac *AuthController) SignOut(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	code := ctx.Param("code")
	verificationCode := utils.Encode(code)

	var user models.User
	result := ac.DB.First(&user, "verification_code = ?", verificationCode)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification code"})
	}

	if user.Verified {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email already verified"})
	}

	user.VerificationCode = ""
	user.Verified = true
	ac.DB.Save(user)

	ctx.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}
