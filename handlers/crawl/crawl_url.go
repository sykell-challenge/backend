package crawl

import (
	"fmt"

	"sykell-challenge/backend/utils"

	"github.com/gin-gonic/gin"
)

// HandleCrawlURL starts a new crawl job or returns existing data
func HandleCrawlURL(g *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	url := req.URL
	fmt.Println("Starting crawl job for URL: ", url)

	// Ping the URL to check if it's available before starting crawl
	fmt.Println("Checking URL availability: ", url)
	pingResult := utils.PingURL(url)

	if !pingResult.Available {
		errorMsg := fmt.Sprintf("URL is not available: %s", pingResult.Error)
		fmt.Println("URL ping failed: ", errorMsg)
		g.JSON(400, gin.H{
			"error": errorMsg,
			"details": gin.H{
				"url":           url,
				"status_code":   pingResult.StatusCode,
				"response_time": pingResult.ResponseTime.Milliseconds(),
				"final_url":     pingResult.FinalURL,
			},
		})
		return
	}

	fmt.Printf("URL is available (Status: %d, Response time: %v)\n", pingResult.StatusCode, pingResult.ResponseTime)

	crawlManager := getCrawlManagerOrError(g)
	if crawlManager == nil {
		return
	}

	// Start the crawl job or get existing data
	job, err := crawlManager.StartCrawlJob(url)
	if err != nil {
		g.JSON(500, gin.H{"error": fmt.Sprintf("Failed to start crawl job: %v", err)})
		return
	}

	// Check if this is returning existing data (mock job with immediate result)
	select {
	case result := <-job.Result:
		if result.Success && job.Status == "completed" {
			// This is existing data, return it immediately
			data := gin.H{
				"message":     "URL already exists in database, returning existing data",
				"is_existing": true,
				"crawl_jobId": job.ID,
				"url_id":      job.URLID,
				"status":      "completed",
				"data": gin.H{
					"title":        result.Data.Title,
					"status_code":  result.Data.StatusCode,
					"html_version": result.Data.HTMLVersion,
					"login_form":   result.Data.LoginForm,
					"tags":         result.Data.Tags,
					"links":        result.Data.Links,
					"last_crawled": result.Data.UpdatedAt,
				},
			}
			g.JSON(200, data)
			return
		}
		// If we get here, there was an error with existing data
		g.JSON(500, gin.H{"error": "Error retrieving existing data"})
		return
	default:
		// This is a new crawl job, return the job info
		data := gin.H{
			"message":    "Crawling started in background",
			"isExisting": false,
			"jobId":      job.ID,
			"urlId":      job.URLID,
			"status":     job.Status,
			"startTime":  job.StartTime,
		}
		g.JSON(200, data)
	}
}
