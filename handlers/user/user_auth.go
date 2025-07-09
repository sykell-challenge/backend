package user

import (
	"net/http"
	"sykell-challenge/backend/auth"
	"sykell-challenge/backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginUser handles POST /users/login
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Provide human-readable error messages
		if req.Username == "" && req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		} else if req.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		} else if req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		}
		return
	}

	// Get user by username
	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update last login timestamp
	if err := h.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log error but don't fail the login
		// In a real application, you might want to log this properly
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user.ToResponse(),
		"token":   token,
	})
}
