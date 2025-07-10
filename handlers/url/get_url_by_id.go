package url

import (
	"sykell-challenge/backend/helpers"

	"github.com/gin-gonic/gin"
)

// GET /urls/:id - Get URL by ID
func (h *URLHandler) GetURLByID(c *gin.Context) {
	id, ok := helpers.ParseIDParam(c, "id")
	if !ok {
		return
	}

	url, err := h.urlRepo.GetByID(id)
	if helpers.HandleDBError(c, err, "URL not found") {
		return
	}

	helpers.SendSuccessResponse(c, gin.H{"data": url})
}
