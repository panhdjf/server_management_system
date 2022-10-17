package models

import (
	"time"
	// "github.com/google/uuid"
)

type User struct {
	ID        int    `gorm:"type:serial;primary_key"`
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SignUpInput struct { //register
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct { //login
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
