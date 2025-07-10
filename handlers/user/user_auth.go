package user

import (
	"net/http"
	"sykell-challenge/backend/helpers"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/validators"

	"github.com/gin-gonic/gin"
)

// LoginUser handles POST /users/login
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Validate request
	if err := validators.ValidateUserLoginRequest(req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Authenticate user
	user, token, err := h.authService.AuthenticateUser(req)
	if err != nil {
		helpers.HandleError(c, err, "Authentication failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user.ToResponse(),
		"token":   token,
	})
}
