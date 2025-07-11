package crawl

import (
	"time"

	"sykell-challenge/backend/models"
)

type Tag struct {
	TagName string `json:"tagName"`
	Count   int    `json:"count"`
}

type Stats struct {
	NumberOfInternalLinks int `json:"numberOfInternalLinks"`
	NumberOfExternalLinks int `json:"numberOfExternalLinks"`
	NumberOfBrokenLinks   int `json:"numberOfBrokenLinks"`
}

type CrawlResponse struct {
	InternalLinks []string `json:"internalLinks"`
	ExternalLinks []string `json:"externalLinks"`
	BrokenLinks   []string `json:"brokenLinks"`
	Title         string   `json:"title"`
	Stats         Stats    `json:"stats"`
	Tags          []Tag    `json:"tags"`
}

// CrawlJob represents a crawl job with its channel and metadata
type CrawlJob struct {
	ID        string
	URL       string
	URLID     uint // Database ID of the URL record
	Status    string
	StartTime time.Time
	Done      chan bool
	Cancel    chan bool
	Result    chan CrawlResult
}

// CrawlResult holds the result of a crawl operation
type CrawlResult struct {
	Success bool
	Error   string
	Data    models.URL
}

// CrawlData represents data from the crawl process
type CrawlData struct {
	MainData  models.URL
	LinkCount int
}

// CrawlJobOrExistingData represents either a new crawl job or existing URL data
type CrawlJobOrExistingData struct {
	Job          *CrawlJob   `json:"job,omitempty"`
	ExistingData *models.URL `json:"existing_data,omitempty"`
	IsExisting   bool        `json:"is_existing"`
}
