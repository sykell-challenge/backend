package user

import (
	"sykell-challenge/backend/repositories"

	"gorm.io/gorm"
)

type UserHandler struct {
	userRepo *repositories.UserRepository
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		userRepo: repositories.NewUserRepository(db),
	}
}
