package crawl

import (
	"context"
	"fmt"
	"net/http"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/services/crawl"
	"sykell-challenge/backend/services/taskq"

	"github.com/gin-gonic/gin"
)

type CrawlRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *CrawlHandler) HandleCrawlURL(g *gin.Context) {
	fmt.Println("HandleCrawlURL called")
	var request CrawlRequest

	if err := g.ShouldBindJSON(&request); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingURL, err := h.urlRepo.GetByURL(request.URL)
	if err == nil && existingURL != nil {
		g.JSON(http.StatusOK, gin.H{
			"message": "URL already crawled",
			"data":    existingURL,
		})
		return
	}

	newURL := models.URL{
		URL:    request.URL,
		Status: "queued",
	}

	if err := h.urlRepo.Create(&newURL); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create URL record"})
		return
	}

	crawlTask := crawl.CreateCrawlTask(request.URL, newURL.ID)

	newURL.CrawlJobID = crawlTask.GetJobID()
	if err := h.urlRepo.Update(&newURL); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL with job ID"})
		return
	}

	ctx := context.Background()
	taskID, err := taskq.EnqueueTask(ctx, crawlTask)
	if err != nil {

		h.urlRepo.UpdateStatus(newURL.ID, "error")
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue crawl task"})
		return
	}

	g.JSON(http.StatusCreated, gin.H{
		"message": "Crawl job queued successfully",
		"jobId":   crawlTask.GetJobID(),
		"taskId":  taskID,
		"urlId":   newURL.ID,
		"url":     request.URL,
		"status":  "queued",
	})
}
