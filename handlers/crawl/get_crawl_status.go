package crawl

import (
	"github.com/gin-gonic/gin"
)

func HandleGetCrawlStatus(g *gin.Context) {
	jobID := getJobIDParam(g)

	crawlManager := getCrawlManagerOrError(g)
	if crawlManager == nil {
		return
	}

	response := gin.H{
		"jobId": jobID,
	}

	// Get crawl job information from database (includes timing info)
	crawlJobRecord, err := crawlManager.GetCrawlJobRepo().GetByJobID(jobID)
	if err != nil {
		sendNotFoundError(g, "Crawl job not found")
		return
	}

	// Add timing information
	response["created_at"] = crawlJobRecord.CreatedAt
	response["started_at"] = crawlJobRecord.StartedAt
	response["completed_at"] = crawlJobRecord.CompletedAt
	response["status"] = crawlJobRecord.Status
	response["duration_ms"] = crawlJobRecord.Duration
	response["running_duration_ms"] = crawlJobRecord.GetRunningDuration()
	response["is_active"] = crawlJobRecord.IsActive()

	if crawlJobRecord.ErrorMsg != "" {
		response["error_message"] = crawlJobRecord.ErrorMsg
	}

	// Try to get URL information from database
	if urlRecord, err := crawlManager.GetURLRepo().GetByCrawlJobID(jobID); err == nil {
		response["url"] = urlRecord.URL
		response["url_id"] = urlRecord.ID
		response["crawl_data"] = gin.H{
			"title":       urlRecord.Title,
			"status_code": urlRecord.StatusCode,
			"login_form":  urlRecord.LoginForm,
			"tags":        urlRecord.Tags,
			"links":       urlRecord.Links,
		}
	}

	// Check if job is still active in memory
	job, exists := crawlManager.GetJob(jobID)
	if exists {
		response["in_memory"] = true
		response["memory_start_time"] = job.StartTime

		// Check if job is completed and has results
		select {
		case result := <-job.Result:
			response["completed"] = true
			response["success"] = result.Success
			if result.Success {
				response["result_data"] = result.Data
			} else {
				response["result_error"] = result.Error
			}
		default:
			response["completed"] = false
		}
	} else {
		response["in_memory"] = false
		response["completed"] = !crawlJobRecord.IsActive()
	}

	g.JSON(200, response)
}
