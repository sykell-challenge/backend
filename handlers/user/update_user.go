package user

import (
	"net/http"
	"strconv"
	"sykell-challenge/backend/helpers"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/validators"

	"github.com/gin-gonic/gin"
)

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Validate request
	if err := validators.ValidateUserUpdateRequest(req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		helpers.HandleError(c, err, "Failed to update user")
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
