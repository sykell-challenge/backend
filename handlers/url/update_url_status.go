package url

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
