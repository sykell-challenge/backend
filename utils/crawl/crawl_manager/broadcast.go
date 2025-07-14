package crawl_manager

import (
	"fmt"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/services/socket"
	crawlUtils "sykell-challenge/backend/utils/crawl"
)

type SocketMessage struct {
	JobID       string       `json:"jobId"`
	URL         string       `json:"url"`
	URLID       string       `json:"urlId"`
	Status      string       `json:"status"`
	StartedAt   string       `json:"startedAt,omitempty"`
	CompletedAt string       `json:"completedAt,omitempty"`
	Progress    int          `json:"progress,omitempty"`
	Title       string       `json:"title,omitempty"`
	StatusCode  int          `json:"statusCode,omitempty"`
	HTMLVersion string       `json:"htmlVersion,omitempty"`
	LoginForm   bool         `json:"loginForm,omitempty"`
	LinksCount  int          `json:"linksCount,omitempty"`
	TagsCount   int          `json:"tagsCount,omitempty"`
	Tags        models.Tags  `json:"tags,omitempty"`
	Links       models.Links `json:"links,omitempty"`
	Error       string       `json:"error,omitempty"`
}

func BroadcastJobQueued(jobID, url string, urlID uint) {
	socket.BroadcastCrawlUpdate("crawl_queued", SocketMessage{
		JobID:  jobID,
		URL:    url,
		Status: "queued",
		URLID:  fmt.Sprintf("%d", urlID),
	})
}

func BroadcastJobStarted(job models.CrawlJob) {
	socket.BroadcastCrawlUpdate("crawl_started", SocketMessage{
		JobID:     fmt.Sprintf("%d", job.Model.ID),
		URL:       job.URL,
		URLID:     fmt.Sprintf("%d", job.URLID),
		Status:    job.Status,
		StartedAt: job.StartedAt.Format("2006-01-02 15:04:05"),
		Progress:  job.Progress,
	})
}

func BroadcastHalfCompleted(job models.CrawlJob, crawlData crawlUtils.CrawlData) {
	fmt.Printf("Broadcasting half-completed job: %s", job.ID)
	fmt.Printf("Job details: %+v", job)
	fmt.Printf("Crawl data: %+v", crawlData)
	socket.BroadcastCrawlUpdate("crawl_half_completed", SocketMessage{
		JobID:       fmt.Sprintf("%d", job.ID),
		URL:         job.URL,
		URLID:       fmt.Sprintf("%d", job.URLID),
		Status:      job.Status,
		StartedAt:   job.StartedAt.Format("2006-01-02 15:04:05"),
		Progress:    job.Progress,
		Title:       crawlData.MainData.Title,
		StatusCode:  crawlData.MainData.StatusCode,
		HTMLVersion: crawlData.MainData.HTMLVersion,
		LoginForm:   crawlData.MainData.LoginForm,
		TagsCount:   len(crawlData.MainData.Tags),
		Tags:        crawlData.MainData.Tags,
	})
}

func BroadcastCompleted(job models.CrawlJob, crawlData crawlUtils.CrawlData) {
	socket.BroadcastCrawlUpdate("crawl_completed", SocketMessage{
		JobID:       fmt.Sprintf("%d", job.ID),
		URL:         job.URL,
		URLID:       fmt.Sprintf("%d", job.URLID),
		Status:      job.Status,
		StartedAt:   job.StartedAt.Format("2006-01-02 15:04:05"),
		CompletedAt: job.CompletedAt.Format("2006-01-02 15:04:05"),
		Progress:    job.Progress,
		Title:       crawlData.MainData.Title,
		StatusCode:  crawlData.MainData.StatusCode,
		HTMLVersion: crawlData.MainData.HTMLVersion,
		LoginForm:   crawlData.MainData.LoginForm,
		LinksCount:  crawlData.LinkCount,
		TagsCount:   len(crawlData.MainData.Tags),
		Tags:        crawlData.MainData.Tags,
		Links:       crawlData.MainData.Links,
	})
}

func BroadcastError(job models.CrawlJob, errorMsg string) {
	socket.BroadcastCrawlUpdate("crawl_error", SocketMessage{
		JobID:  fmt.Sprintf("%d", job.ID),
		URL:    job.URL,
		URLID:  fmt.Sprintf("%d", job.URLID),
		Error:  errorMsg,
		Status: "error",
	})
}

func BroadcastCancelled(job models.CrawlJob) {
	socket.BroadcastCrawlUpdate("crawl_cancelled", SocketMessage{
		JobID:  fmt.Sprintf("%d", job.ID),
		URL:    job.URL,
		URLID:  fmt.Sprintf("%d", job.URLID),
		Status: "cancelled",
	})
}
