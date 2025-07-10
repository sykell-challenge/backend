package services

import (
	"errors"
	"sykell-challenge/backend/models"
)

// UpdateUser handles user update business logic
func (s *UserService) UpdateUser(id uint, req models.UserUpdateRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		// Check if new email already exists
		exists, err := s.userRepo.EmailExists(*req.Email)
		if err != nil {
			return nil, err
		}
		if exists && *req.Email != user.Email {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
