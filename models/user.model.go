package models

import (
	"time"
)

type User struct {
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PasswordResetAt   time.Time
	Provider          string `gorm:"not null"`
	Password          string `gorm:"not null"`
	Role              string `gorm:"not null"`
	Email             string `gorm:"uniqueIndex; not null"`
	Photo             string `gorm:"not null"`
	VerificationCode  string
	PasswordResetCode string
	Name              string `gorm:"type:varchar(255);not null"`
	ID                uint   `gorm:"primary_key"`
	Verified          bool   `gorm:"not null"`
}

type SignUp struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,min=8"`
	Photo                string `json:"photo" binding:"required"`
}

type SignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Photo     string    `json:"photo"`
	Provider  string    `json:"provider"`
	ID        uint      `json:"id"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,min=8"`
}
