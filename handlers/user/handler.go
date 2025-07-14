package user

import (
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/repositories"
	"sykell-challenge/backend/services"
)

type UserHandler struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserHandler() *UserHandler {
	db := db.GetDB()
	userRepo := repositories.NewUserRepository(db)
	return &UserHandler{
		userService: services.NewUserService(userRepo),
		authService: services.NewAuthService(userRepo),
	}
}
