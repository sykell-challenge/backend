package crawl_manager

import (
	"sykell-challenge/backend/services/socket"
	crawlUtils "sykell-challenge/backend/utils/crawl"
)

func BroadcastJobQueued(jobID, url string, urlID uint) {
	socket.BroadcastCrawlUpdate("crawl_queued", map[string]interface{}{
		"jobId":  jobID,
		"url":    url,
		"status": "queued",
		"url_id": urlID,
	})
}

func BroadcastJobStarted(job crawlUtils.CrawlJob) {
	socket.BroadcastCrawlUpdate("crawl_started", map[string]interface{}{
		"jobId":  job.ID,
		"url":    job.URL,
		"status": "running",
		"url_id": job.URLID,
	})
}

func BroadcastHalfCompleted(job crawlUtils.CrawlJob, crawlData crawlUtils.CrawlData) {
	socket.BroadcastCrawlUpdate("crawl_half_completed", map[string]interface{}{
		"jobId":        job.ID,
		"url":          job.URL,
		"url_id":       job.URLID,
		"title":        crawlData.MainData.Title,
		"status_code":  crawlData.MainData.StatusCode,
		"html_version": crawlData.MainData.HTMLVersion,
		"login_form":   crawlData.MainData.LoginForm,
		"tags_count":   len(crawlData.MainData.Tags),
		"links_count":  crawlData.LinkCount,
		"status":       "half_completed",
	})
}

func BroadcastCompleted(job crawlUtils.CrawlJob, crawlData crawlUtils.CrawlData) {
	socket.BroadcastCrawlUpdate("crawl_completed", map[string]interface{}{
		"jobId":        job.ID,
		"url":          job.URL,
		"url_id":       job.URLID,
		"title":        crawlData.MainData.Title,
		"status_code":  crawlData.MainData.StatusCode,
		"html_version": crawlData.MainData.HTMLVersion,
		"login_form":   crawlData.MainData.LoginForm,
		"tags_count":   len(crawlData.MainData.Tags),
		"links_count":  crawlData.LinkCount,
		"status":       "completed",
		"links":        crawlData.MainData.Links,
	})
}

func BroadcastError(job crawlUtils.CrawlJob, errorMsg string) {
	socket.BroadcastCrawlUpdate("crawl_error", map[string]interface{}{
		"jobId":  job.ID,
		"url":    job.URL,
		"url_id": job.URLID,
		"error":  errorMsg,
		"status": "error",
	})
}

func BroadcastCancelled(job crawlUtils.CrawlJob) {
	socket.BroadcastCrawlUpdate("crawl_cancelled", map[string]interface{}{
		"jobId":  job.ID,
		"url":    job.URL,
		"url_id": job.URLID,
		"status": "cancelled",
	})
}
