package user

import (
	"net/http"
	"strconv"
	"sykell-challenge/backend/helpers"

	"github.com/gin-gonic/gin"
)

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Delete user
	if err := h.userService.DeleteUser(uint(id)); err != nil {
		helpers.HandleError(c, err, "Failed to delete user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
