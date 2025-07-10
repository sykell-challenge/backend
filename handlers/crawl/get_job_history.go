package crawl

import (
	"github.com/gin-gonic/gin"
)

func HandleGetJobHistory(g *gin.Context) {
	crawlManager := getCrawlManagerOrError(g)
	if crawlManager == nil {
		return
	}

	// Get limit from query parameter, default to 50, max 1000
	limit := parseLimitQuery(g, 50, 1000)

	jobs, err := crawlManager.GetCrawlJobRepo().GetJobHistory(limit)
	if err != nil {
		sendInternalError(g, "Failed to fetch job history")
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
