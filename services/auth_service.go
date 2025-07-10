package services

import (
	"errors"
	"sykell-challenge/backend/auth"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// AuthenticateUser handles user authentication business logic
func (s *AuthService) AuthenticateUser(req models.UserLoginRequest) (*models.User, string, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", errors.New("account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Update last login timestamp
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log error but don't fail the login
		// In a real application, you might want to log this properly
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}
