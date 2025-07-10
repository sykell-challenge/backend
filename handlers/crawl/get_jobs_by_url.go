package crawl

import (
	"github.com/gin-gonic/gin"
)

// HandleGetJobsByURL returns all crawl jobs for a specific URL ID
func HandleGetJobsByURL(g *gin.Context) {
	crawlManager := getCrawlManagerOrError(g)
	if crawlManager == nil {
		return
	}

	urlID, ok := parseURLIDParam(g)
	if !ok {
		return
	}

	jobs, err := crawlManager.GetCrawlJobRepo().GetJobsByURLID(urlID)
	if err != nil {
		sendInternalError(g, "Failed to fetch jobs for URL")
		return
	}

	jobList := make([]gin.H, 0, len(jobs))
	for _, job := range jobs {
		jobInfo := buildJobInfoFromCrawlJob(&job)
		jobList = append(jobList, jobInfo)
	}

	g.JSON(200, gin.H{
		"jobs":   jobList,
		"count":  len(jobList),
		"url_id": urlID,
	})
}
