package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"uniqueIndex; not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null"`
	Provider  string `gorm:"not null"`
	Photo     string `gorm:"not null"`
	Verified  bool   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Photo     string    `json:"photo"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
