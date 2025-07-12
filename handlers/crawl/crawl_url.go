package crawl

import (
	"context"
	"log"
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
	var request CrawlRequest
	if err := h.validateRequeest(g, &request); err != nil {
		return
	}
	log.Printf("request url: %v", request)

	if crawled := h.isUrlAlreadyCrawled(g, request.URL); crawled {
		return
	}

	newURL := models.URL{
		URL:    request.URL,
		Status: "queued",
	}

	// store newURL in database
	if err := h.urlRepo.Create(&newURL); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create URL record"})
		return
	}
	g.JSON(200, newURL)

	crawlTask := crawl.CreateCrawlTask(request.URL, newURL.ID)

	// update url in database with jobid
	newURL.JobId = crawlTask.GetJobID()
	if err := h.urlRepo.Update(&newURL); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL with job ID"})
		return
	}

	// start crawling in background
	ctx := context.Background()
	taskID, err := taskq.EnqueueTask(ctx, crawlTask)
	if err != nil {
		h.urlRepo.UpdateStatus(newURL.ID, "error")
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue crawl task"})
		return
	}

	// send response
	g.JSON(http.StatusCreated, gin.H{
		"message": "Crawl job queued successfully",
		"jobId":   crawlTask.GetJobID(),
		"taskId":  taskID,
		"urlId":   newURL.ID,
		"url":     request.URL,
		"status":  "queued",
	})
}

func (h *CrawlHandler) validateRequeest(g *gin.Context, request *CrawlRequest) error {

	if err := g.ShouldBindJSON(&request); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func (h *CrawlHandler) isUrlAlreadyCrawled(g *gin.Context, urlString string) bool {

	existingURL, err := h.urlRepo.GetByURL(urlString)
	if err == nil && existingURL != nil {

		g.JSON(http.StatusOK, existingURL)
		return true
	}

	return false
}
