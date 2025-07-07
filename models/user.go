package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string     `json:"username" gorm:"type:varchar(255);uniqueIndex;not null"`
	Email        string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string     `json:"-" gorm:"type:varchar(255);not null"` // "-" excludes from JSON serialization
	FirstName    string     `json:"first_name" gorm:"type:varchar(100)"`
	LastName     string     `json:"last_name" gorm:"type:varchar(100)"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}

// UserResponse represents the user data returned in API responses (excludes sensitive fields)
type UserResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToResponse converts a User model to UserResponse (excludes sensitive fields)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// UserCreateRequest represents the data needed to create a new user
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"max=100"`
	LastName  string `json:"last_name" binding:"max=100"`
}

// UserUpdateRequest represents the data that can be updated for a user
type UserUpdateRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=100"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=100"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UserLoginRequest represents the data needed for user login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
