package crawl

import (
	"fmt"
	"strconv"

	"sykell-challenge/backend/utils"
	"sykell-challenge/backend/utils/crawl"

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

	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
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

// HandleGetCrawlStatus returns the status of a crawl job
func HandleGetCrawlStatus(g *gin.Context) {
	jobID := g.Param("jobId")

	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
		return
	}

	response := gin.H{
		"jobId": jobID,
	}

	// Get crawl job information from database (includes timing info)
	crawlJobRecord, err := crawlManager.GetCrawlJobRepo().GetByJobID(jobID)
	if err != nil {
		g.JSON(404, gin.H{"error": "Crawl job not found"})
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

// HandleCancelCrawl cancels a running crawl job
func HandleCancelCrawl(g *gin.Context) {
	jobID := g.Param("jobId")

	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
		return
	}

	success := crawlManager.CancelJob(jobID)
	if !success {
		// Check if job exists in database but not in memory (already completed)
		if crawlJobRecord, err := crawlManager.GetCrawlJobRepo().GetByJobID(jobID); err == nil {
			if crawlJobRecord.IsActive() {
				// Update status to cancelled in database
				crawlJobRecord.SetCancelled()
				crawlManager.GetCrawlJobRepo().Update(crawlJobRecord)

				// Also update URL status
				if urlRecord, err := crawlManager.GetURLRepo().GetByCrawlJobID(jobID); err == nil {
					crawlManager.GetURLRepo().UpdateStatus(urlRecord.ID, "error")
				}
			}
			g.JSON(200, gin.H{"message": "Crawl job was already completed, status updated in database"})
			return
		}
		g.JSON(404, gin.H{"error": "Crawl job not found"})
		return
	}

	g.JSON(200, gin.H{"message": "Crawl job cancelled successfully"})
}

// HandleListActiveJobs returns all currently active crawl jobs
func HandleListActiveJobs(g *gin.Context) {
	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
		return
	}

	// Get active jobs from database (includes timing info)
	activeJobsFromDB, err := crawlManager.GetCrawlJobRepo().GetActiveJobs()
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to fetch active jobs from database"})
		return
	}

	// Get in-memory jobs for comparison
	memoryJobs := crawlManager.ListActiveJobs()

	jobList := make([]gin.H, 0, len(activeJobsFromDB))
	for _, dbJob := range activeJobsFromDB {
		jobInfo := gin.H{
			"jobId":               dbJob.JobID,
			"url":                 dbJob.URL,
			"url_id":              dbJob.URLID,
			"status":              dbJob.Status,
			"created_at":          dbJob.CreatedAt,
			"started_at":          dbJob.StartedAt,
			"running_duration_ms": dbJob.GetRunningDuration(),
			"is_active":           dbJob.IsActive(),
		}

		// Check if job is also in memory
		if memJob, exists := memoryJobs[dbJob.JobID]; exists {
			jobInfo["in_memory"] = true
			jobInfo["memory_start_time"] = memJob.StartTime
			jobInfo["memory_status"] = memJob.Status
		} else {
			jobInfo["in_memory"] = false
		}

		jobList = append(jobList, jobInfo)
	}

	g.JSON(200, gin.H{
		"active_jobs":   jobList,
		"count":         len(jobList),
		"memory_jobs":   len(memoryJobs),
		"database_jobs": len(activeJobsFromDB),
	})
}

// HandleGetJobHistory returns historical crawl jobs
func HandleGetJobHistory(g *gin.Context) {
	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
		return
	}

	// Get limit from query parameter, default to 50
	limit := 50
	if limitStr := g.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	jobs, err := crawlManager.GetCrawlJobRepo().GetJobHistory(limit)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to fetch job history"})
		return
	}

	jobList := make([]gin.H, 0, len(jobs))
	for _, job := range jobs {
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

		// Try to get URL data for this job
		if urlRecord, err := crawlManager.GetURLRepo().GetByCrawlJobID(job.JobID); err == nil {
			jobInfo["data"] = gin.H{
				"title":        urlRecord.Title,
				"status_code":  urlRecord.StatusCode,
				"html_version": urlRecord.HTMLVersion,
				"login_form":   urlRecord.LoginForm,
				"tags_count":   len(urlRecord.Tags),
				"links_count":  len(urlRecord.Links),
				"url_status":   urlRecord.Status,
			}

			// Include full data by default, exclude only if explicitly set to false
			includeFullData := g.Query("include_full_data") != "false"
			if includeFullData {
				jobInfo["data"].(gin.H)["tags"] = urlRecord.Tags
				jobInfo["data"].(gin.H)["links"] = urlRecord.Links
			}
		} else {
			// If no URL data found, indicate it
			jobInfo["data"] = nil
		}

		jobList = append(jobList, jobInfo)
	}

	g.JSON(200, gin.H{
		"job_history":       jobList,
		"count":             len(jobList),
		"limit":             limit,
		"include_full_data": g.Query("include_full_data") != "false",
		"note":              "Full tags and links data included by default. Use ?include_full_data=false to get summary only",
	})
}

// HandleGetJobsByURL returns all crawl jobs for a specific URL ID
func HandleGetJobsByURL(g *gin.Context) {
	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
		return
	}

	urlIDStr := g.Param("urlId")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		g.JSON(400, gin.H{"error": "Invalid URL ID format"})
		return
	}

	jobs, err := crawlManager.GetCrawlJobRepo().GetJobsByURLID(uint(urlID))
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to fetch jobs for URL"})
		return
	}

	jobList := make([]gin.H, 0, len(jobs))
	for _, job := range jobs {
		jobInfo := gin.H{
			"jobId":               job.JobID,
			"url":                 job.URL,
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

		jobList = append(jobList, jobInfo)
	}

	g.JSON(200, gin.H{
		"jobs":   jobList,
		"count":  len(jobList),
		"url_id": uint(urlID),
	})
}

// HandleGetJobStats returns crawl job statistics
func HandleGetJobStats(g *gin.Context) {
	stats, err := crawl.GetJobStats()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"stats": stats})
}

// HandleGetURLByJobID returns URL data for a specific job ID
func HandleGetURLByJobID(g *gin.Context) {
	jobID := g.Param("jobId")

	crawlManager := crawl.GetCrawlManager()
	if crawlManager == nil {
		g.JSON(500, gin.H{"error": "Crawl manager not initialized"})
		return
	}

	// First check if the job exists in the database
	crawlJobRecord, err := crawlManager.GetCrawlJobRepo().GetByJobID(jobID)
	if err != nil {
		g.JSON(404, gin.H{"error": "Crawl job not found"})
		return
	}

	// Get the URL record associated with this job
	urlRecord, err := crawlManager.GetURLRepo().GetByCrawlJobID(jobID)
	if err != nil {
		g.JSON(404, gin.H{"error": "URL data not found for this job"})
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
