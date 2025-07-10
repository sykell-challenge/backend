package user

import (
	"net/http"
	"sykell-challenge/backend/helpers"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/validators"

	"github.com/gin-gonic/gin"
)

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Validate request
	if err := validators.ValidateUserCreateRequest(req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Create user
	user, err := h.userService.CreateUser(req)
	if err != nil {
		helpers.HandleError(c, err, "Failed to create user")
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}
