package url

import (
	"net/http"
	"strconv"
	"sykell-challenge/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PUT /urls/:id - Update URL
func (h *URLHandler) UpdateURL(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Check if URL exists
	existingURL, err := h.urlRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateData models.URL
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status if provided
	if updateData.Status != "" {
		if updateData.Status != "queued" && updateData.Status != "running" && updateData.Status != "done" && updateData.Status != "error" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be one of: queued, running, done, error"})
			return
		}
	}

	// Update fields
	existingURL.ID = uint(id) // Ensure ID is set for update
	if updateData.URL != "" {
		existingURL.URL = updateData.URL
	}
	if updateData.Status != "" {
		existingURL.Status = updateData.Status
	}
	if updateData.HTMLVersion != "" {
		existingURL.HTMLVersion = updateData.HTMLVersion
	}
	// LoginForm can be updated (boolean field)
	existingURL.LoginForm = updateData.LoginForm

	if err := h.urlRepo.Update(existingURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": existingURL})
}
