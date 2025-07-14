package crawl

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/services/crawl"
	"sykell-challenge/backend/services/taskq"
	"sykell-challenge/backend/utils"

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

	if isUrlAvailable, statusCode := utils.IsURLAvailable(request.URL); !isUrlAvailable {
		g.JSON(statusCode, gin.H{"error": gin.H{"message": "URL is not available", "code": statusCode}})
		return
	}

	newURL := models.URL{
		URL:    request.URL,
		Status: "queued",
	}

	// store newURL in database
	if err := h.urlRepo.Create(&newURL); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"message": "Failed to create URL record",
			"code":    http.StatusInternalServerError,
		}})
		return
	}

	crawlTask := crawl.CreateCrawlTask(request.URL, newURL.ID)

	// update url in database with jobid
	newURL.JobId = fmt.Sprintf("%d", crawlTask.CrawlJob.ID)
	if err := h.urlRepo.Update(&newURL); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"message": "Failed to update URL with job ID",
			"code":    http.StatusInternalServerError,
		}})
		return
	}

	// start crawling in background
	ctx := context.Background()
	_, err := taskq.EnqueueTask(ctx, crawlTask)
	if err != nil {
		h.urlRepo.UpdateStatus(newURL.ID, "error")
		g.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"message": "Failed to enqueue crawl task",
			"code":    http.StatusInternalServerError,
		}})
		return
	}

	// send response
	g.JSON(http.StatusCreated, gin.H{
		"alreadyCrawled": false,
		"data":           newURL,
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

		g.JSON(http.StatusOK, gin.H{
			"alreadyCrawled": true,
			"data":           existingURL,
		})
		return true
	}

	return false
}
