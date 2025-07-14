package crawl

import (
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/repositories"

	"gorm.io/gorm"
)

func NewCrawlHandler() *CrawlHandler {
	db := db.GetDB()

	return &CrawlHandler{
		db:      db,
		urlRepo: repositories.NewURLRepository(db),
		jobRepo: repositories.NewCrawlJobRepository(db),
	}
}

type CrawlHandler struct {
	db      *gorm.DB
	urlRepo *repositories.URLRepository
	jobRepo *repositories.CrawlJobRepository
}
