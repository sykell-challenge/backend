package crawl

import (
	"net/http"
	"sykell-challenge/backend/repositories"

	"github.com/gin-gonic/gin"
)

// HandleGetAllCrawlJobs returns all past crawl jobs
func (h *CrawlHandler) HandleGetAllCrawlJobs(c *gin.Context) {
	repo := repositories.NewCrawlJobRepository(h.db)

	jobs, err := repo.GetJobHistory(100) // limit to 100 most recent jobs
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}
