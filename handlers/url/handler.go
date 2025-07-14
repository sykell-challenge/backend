package url

import (
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/repositories"
)

type URLHandler struct {
	urlRepo *repositories.URLRepository
}

func NewURLHandler() *URLHandler {
	db := db.GetDB()
	return &URLHandler{
		urlRepo: repositories.NewURLRepository(db),
	}
}
