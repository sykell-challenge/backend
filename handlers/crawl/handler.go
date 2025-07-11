package crawl

import (
	"sykell-challenge/backend/repositories"

	"gorm.io/gorm"
)

func NewCrawlHandler(db *gorm.DB) *CrawlHandler {
	return &CrawlHandler{
		db:      db,
		urlRepo: repositories.NewURLRepository(db),
	}
}

type CrawlHandler struct {
	db      *gorm.DB
	urlRepo *repositories.URLRepository
}
