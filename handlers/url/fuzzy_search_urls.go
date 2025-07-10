package url

import (
	"sykell-challenge/backend/helpers"

	"github.com/gin-gonic/gin"
)

// GET /urls/search/fuzzy?q=... - Fuzzy search URLs
func (h *URLHandler) FuzzySearchURLs(c *gin.Context) {
	// Using the helper function for query parameter validation
	query, ok := helpers.RequireQueryParam(c, "q")
	if !ok {
		return
	}

	// Using the helper function for limit parsing
	limit := helpers.ParseLimitQuery(c, 10, 100)

	urls, err := h.urlRepo.SearchURLs(query, limit)
	if err != nil {
		// Using the helper function for error response
		helpers.SendInternalError(c, err.Error())
		return
	}

	// Using the helper function for success response
	helpers.SendSuccessResponse(c, gin.H{"data": urls})
}
