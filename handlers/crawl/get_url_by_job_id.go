package crawl

import (
	"github.com/gin-gonic/gin"
)

func HandleGetURLByJobID(g *gin.Context) {
	jobID := getJobIDParam(g)

	crawlManager := getCrawlManagerOrError(g)
	if crawlManager == nil {
		return
	}

	// First check if the job exists in the database
	crawlJobRecord, err := crawlManager.GetCrawlJobRepo().GetByJobID(jobID)
	if err != nil {
		sendNotFoundError(g, "Crawl job not found")
		return
	}

	// Get the URL record associated with this job
	urlRecord, err := crawlManager.GetURLRepo().GetByCrawlJobID(jobID)
	if err != nil {
		sendNotFoundError(g, "URL data not found for this job")
		return
	}

	// Build the response
	response := gin.H{
		"jobId":      jobID,
		"url_id":     urlRecord.ID,
		"url":        urlRecord.URL,
		"status":     urlRecord.Status,
		"job_status": crawlJobRecord.Status,
		"data": gin.H{
			"title":        urlRecord.Title,
			"status_code":  urlRecord.StatusCode,
			"html_version": urlRecord.HTMLVersion,
			"login_form":   urlRecord.LoginForm,
			"tags":         urlRecord.Tags,
			"links":        urlRecord.Links,
		},
		"timestamps": gin.H{
			"created_at":    urlRecord.CreatedAt,
			"updated_at":    urlRecord.UpdatedAt,
			"job_created":   crawlJobRecord.CreatedAt,
			"job_started":   crawlJobRecord.StartedAt,
			"job_completed": crawlJobRecord.CompletedAt,
		},
		"timing": gin.H{
			"duration_ms":         crawlJobRecord.Duration,
			"running_duration_ms": crawlJobRecord.GetRunningDuration(),
			"is_active":           crawlJobRecord.IsActive(),
		},
	}

	// Add error message if present
	if crawlJobRecord.ErrorMsg != "" {
		response["error_message"] = crawlJobRecord.ErrorMsg
	}

	// Check if job is still active in memory
	if job, exists := crawlManager.GetJob(jobID); exists {
		response["in_memory"] = true
		response["memory_status"] = job.Status
		response["memory_start_time"] = job.StartTime
	} else {
		response["in_memory"] = false
	}

	g.JSON(200, response)
}
