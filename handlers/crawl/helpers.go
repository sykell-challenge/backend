package crawl

import (
	"strconv"

	"sykell-challenge/backend/models"
	"sykell-challenge/backend/utils/crawl"

	"github.com/gin-gonic/gin"
)

// getCrawlManagerOrError returns the crawl manager or sends an error response and returns nil
func getCrawlManagerOrError(g *gin.Context) *crawl.CrawlManager {
	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
		return nil
	}
	return crawlManager
}

// parseURLIDParam extracts and validates URL ID from URL parameter
func parseURLIDParam(g *gin.Context) (uint, bool) {
	urlIDStr := g.Param("urlId")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		g.JSON(400, gin.H{"error": "Invalid URL ID format"})
		return 0, false
	}
	return uint(urlID), true
}

// getJobIDParam extracts job ID from URL parameter
func getJobIDParam(g *gin.Context) string {
	return g.Param("jobId")
}

// buildJobInfo creates a common job info structure
func buildJobInfo(job interface{}) gin.H {
	// This is a generic helper - you'd need to use type assertion or interface
	// For now, let's create specific builders for different job types
	return gin.H{}
}

// buildJobInfoFromCrawlJob creates job info from crawl job record
func buildJobInfoFromCrawlJob(job *models.CrawlJob) gin.H {
	jobInfo := gin.H{
		"jobId":               job.JobID,
		"url":                 job.URL,
		"url_id":              job.URLID,
		"status":              job.Status,
		"created_at":          job.CreatedAt,
		"started_at":          job.StartedAt,
		"completed_at":        job.CompletedAt,
		"duration_ms":         job.Duration,
		"running_duration_ms": job.GetRunningDuration(),
		"is_active":           job.IsActive(),
	}

	if job.ErrorMsg != "" {
		jobInfo["error_message"] = job.ErrorMsg
	}

	return jobInfo
}

// parseLimitQuery parses limit query parameter with default and max values
func parseLimitQuery(g *gin.Context, defaultLimit, maxLimit int) int {
	limit := defaultLimit
	if limitStr := g.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= maxLimit {
			limit = parsedLimit
		}
	}
	return limit
}

// sendNotFoundError sends a 404 error response
func sendNotFoundError(g *gin.Context, message string) {
	g.JSON(404, gin.H{"error": message})
}

// sendInternalError sends a 500 error response
func sendInternalError(g *gin.Context, message string) {
	g.JSON(500, gin.H{"error": message})
}

// sendBadRequestError sends a 400 error response
func sendBadRequestError(g *gin.Context, message string) {
	g.JSON(400, gin.H{"error": message})
}

// sendSuccessResponse sends a 200 success response
func sendSuccessResponse(g *gin.Context, data gin.H) {
	g.JSON(200, data)
}

// checkIncludeFullData checks if full data should be included based on query parameter
func checkIncludeFullData(g *gin.Context) bool {
	return g.Query("include_full_data") != "false"
}

// addURLDataToJobInfo adds URL data to job info if available
func addURLDataToJobInfo(crawlManager *crawl.CrawlManager, jobInfo gin.H, jobID string, includeFullData bool) {
	if urlRecord, err := crawlManager.GetURLRepo().GetByCrawlJobID(jobID); err == nil {
		jobInfo["data"] = gin.H{
			"title":        urlRecord.Title,
			"status_code":  urlRecord.StatusCode,
			"html_version": urlRecord.HTMLVersion,
			"login_form":   urlRecord.LoginForm,
			"tags_count":   len(urlRecord.Tags),
			"links_count":  len(urlRecord.Links),
			"url_status":   urlRecord.Status,
		}

		if includeFullData {
			jobInfo["data"].(gin.H)["tags"] = urlRecord.Tags
			jobInfo["data"].(gin.H)["links"] = urlRecord.Links
		}
	} else {
		jobInfo["data"] = nil
	}
}

// addMemoryJobInfo adds memory job information to response
func addMemoryJobInfo(crawlManager *crawl.CrawlManager, response gin.H, jobID string) {
	if job, exists := crawlManager.GetJob(jobID); exists {
		response["in_memory"] = true
		response["memory_status"] = job.Status
		response["memory_start_time"] = job.StartTime
	} else {
		response["in_memory"] = false
	}
}
