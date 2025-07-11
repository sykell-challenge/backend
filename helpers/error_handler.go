package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleError(c *gin.Context, err error, defaultMessage string) {
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Check for specific error messages
	switch err.Error() {
	case "username already exists":
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
	case "email already exists":
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
	case "invalid credentials":
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case "account is deactivated":
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
	case "failed to generate token":
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultMessage})
	}
}

func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
