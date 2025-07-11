package crawl

import (
	"log"
	"net/http"
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/repositories"
	"sykell-challenge/backend/services/socket"
	"sykell-challenge/backend/services/taskq"

	"github.com/gin-gonic/gin"
)

func (h *CrawlHandler) HandleCancelCrawl(g *gin.Context) {
	jobID := g.Param("jobId")

	if jobID == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	// Initialize database and repositories
	database := db.GetDB()
	urlRepo := repositories.NewURLRepository(database)

	// Find the URL record by job ID
	urlRecord, err := urlRepo.GetByJobID(jobID)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Check if the job is cancellable
	if urlRecord.Status == "done" {
		g.JSON(http.StatusConflict, gin.H{"error": "Job already completed"})
		return
	}

	if urlRecord.Status == "cancelled" {
		g.JSON(http.StatusConflict, gin.H{"error": "Job already cancelled"})
		return
	}

	if urlRecord.Status == "error" {
		g.JSON(http.StatusConflict, gin.H{"error": "Job already failed"})
		return
	}

	// Try to cancel the running job
	if urlRecord.Status == "running" {
		if taskq.CancelJob(jobID) {
			// Job was running and successfully cancelled
			log.Printf("Cancelled running job: %s", jobID)
		} else {
			// Job might have just finished or wasn't found in running jobs
			log.Printf("Job not found in running jobs, updating status anyway: %s", jobID)
		}
	}

	// Update URL status to cancelled
	if err := urlRepo.UpdateStatus(urlRecord.ID, "cancelled"); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job status"})
		return
	}

	// Broadcast cancellation
	socket.BroadcastCrawlUpdate("crawl_cancelled", map[string]interface{}{
		"jobId":  jobID,
		"url":    urlRecord.URL,
		"url_id": urlRecord.ID,
		"status": "cancelled",
	})

	g.JSON(http.StatusOK, gin.H{
		"message": "Job cancelled successfully",
		"job_id":  jobID,
		"url_id":  urlRecord.ID,
		"status":  "cancelled",
	})
}
