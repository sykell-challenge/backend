package url

import (
	"net/http"
	"strconv"
	"sykell-challenge/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /urls/:id - Get URL by ID
func (h *URLHandler) GetURLByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	url, err := h.urlRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": url})
}

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

// PATCH /urls/:id/status - Update only URL status
func (h *URLHandler) UpdateURLStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	if statusUpdate.Status != "queued" && statusUpdate.Status != "running" && statusUpdate.Status != "done" && statusUpdate.Status != "error" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be one of: queued, running, done, error"})
		return
	}

	if err := h.urlRepo.UpdateStatus(uint(id), statusUpdate.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// DELETE /urls/:id - Delete URL
func (h *URLHandler) DeleteURL(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Check if URL exists
	_, err = h.urlRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.urlRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}
