package services

import (
	"errors"
	"sykell-challenge/backend/models"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles user creation business logic
func (s *UserService) CreateUser(req models.UserCreateRequest) (*models.User, error) {
	// Check if username already exists
	exists, err := s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
