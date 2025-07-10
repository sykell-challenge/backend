package url

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /urls/stats - Get URL statistics
func (h *URLHandler) GetURLStats(c *gin.Context) {
	count, err := h.urlRepo.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get counts by status
	statuses := []string{"queued", "running", "done", "error"}
	statusCounts := make(map[string]int)

	for _, status := range statuses {
		urls, err := h.urlRepo.GetByStatus(status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		statusCounts[status] = len(urls)
	}

	stats := gin.H{
		"total_urls":    count,
		"status_counts": statusCounts,
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
