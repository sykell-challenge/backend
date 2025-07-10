package user

import (
	"net/http"
	"strconv"
	"sykell-challenge/backend/helpers"

	"github.com/gin-gonic/gin"
)

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		helpers.HandleError(c, err, "Failed to fetch user")
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
