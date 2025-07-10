package user

import (
	"sykell-challenge/backend/repositories"
	"sykell-challenge/backend/services"

	"gorm.io/gorm"
)

type UserHandler struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	userRepo := repositories.NewUserRepository(db)
	return &UserHandler{
		userService: services.NewUserService(userRepo),
		authService: services.NewAuthService(userRepo),
	}
}
