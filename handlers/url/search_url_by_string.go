package url

import (
	"sykell-challenge/backend/helpers"

	"github.com/gin-gonic/gin"
)

func (h *URLHandler) SearchURLByString(c *gin.Context) {

	urlString, ok := helpers.RequireQueryParam(c, "url")
	if !ok {
		return
	}

	url, err := h.urlRepo.GetByURL(urlString)

	if helpers.HandleDBError(c, err, "URL not found") {
		return
	}

	helpers.SendSuccessResponse(c, gin.H{"data": url})
}
