package url

import (
	"net/http"
	"sykell-challenge/backend/repositories"

	"github.com/gin-gonic/gin"
)

// GET /urls - Get all URLs with pagination, sorting, and filtering
func (h *URLHandler) GetURLs(c *gin.Context) {
	var params repositories.URLQueryParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate pagination parameters
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	result, err := h.urlRepo.GetAllWithParams(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
