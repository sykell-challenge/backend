package validators

import (
	"errors"
	"sykell-challenge/backend/models"
)

// ValidateUserLoginRequest validates login request
func ValidateUserLoginRequest(req models.UserLoginRequest) error {
	if req.Username == "" && req.Password == "" {
		return errors.New("username and password are required")
	}
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// ValidateUserCreateRequest validates user creation request
func ValidateUserCreateRequest(req models.UserCreateRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if req.LastName == "" {
		return errors.New("last name is required")
	}
	return nil
}

// ValidateUserUpdateRequest validates user update request
func ValidateUserUpdateRequest(req models.UserUpdateRequest) error {
	// For update requests, we can allow partial updates
	// so we don't need to validate all fields
	if req.Email != nil && *req.Email == "" {
		return errors.New("email cannot be empty")
	}
	if req.FirstName != nil && *req.FirstName == "" {
		return errors.New("first name cannot be empty")
	}
	if req.LastName != nil && *req.LastName == "" {
		return errors.New("last name cannot be empty")
	}
	return nil
}
