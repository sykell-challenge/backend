package url

import (
	"net/http"
	"sykell-challenge/backend/models"

	"github.com/gin-gonic/gin"
)

// POST /urls - Create new URL
func (h *URLHandler) CreateURL(c *gin.Context) {
	var url models.URL

	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if url.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// Set default status if not provided
	if url.Status == "" {
		url.Status = "queued"
	}

	// Validate status enum
	if url.Status != "queued" && url.Status != "running" && url.Status != "done" && url.Status != "error" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be one of: queued, running, done, error"})
		return
	}

	if err := h.urlRepo.Create(&url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": url})
}
