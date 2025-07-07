package url

import (
	"sykell-challenge/backend/repositories"

	"gorm.io/gorm"
)

type URLHandler struct {
	urlRepo *repositories.URLRepository
}

func NewURLHandler(db *gorm.DB) *URLHandler {
	return &URLHandler{
		urlRepo: repositories.NewURLRepository(db),
	}
}
