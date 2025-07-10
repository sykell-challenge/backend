package url

import (
	"sykell-challenge/backend/helpers"

	"github.com/gin-gonic/gin"
)

// GET /urls/search?url=... - Search URL by URL string
func (h *URLHandler) SearchURLByString(c *gin.Context) {
	// Using the helper function for query parameter validation
	urlString, ok := helpers.RequireQueryParam(c, "url")
	if !ok {
		return
	}

	url, err := h.urlRepo.GetByURL(urlString)
	// Using the helper function for database error handling
	if helpers.HandleDBError(c, err, "URL not found") {
		return
	}

	// Using the helper function for success response
	helpers.SendSuccessResponse(c, gin.H{"data": url})
}
