package crawl

import (
	"github.com/gin-gonic/gin"
)

func HandleListActiveJobs(g *gin.Context) {
	crawlManager := getCrawlManagerOrError(g)
	if crawlManager == nil {
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
