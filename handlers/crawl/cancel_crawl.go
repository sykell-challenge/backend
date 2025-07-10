package crawl

import (
	"github.com/gin-gonic/gin"
)

func HandleCancelCrawl(g *gin.Context) {
	jobID := getJobIDParam(g)

	crawlManager := getCrawlManagerOrError(g)
	if crawlManager == nil {
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
		sendNotFoundError(g, "Crawl job not found")
		return
	}

	g.JSON(200, gin.H{"message": "Crawl job cancelled successfully"})
}
