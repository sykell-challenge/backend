package crawl

import (
	"fmt"
	"sync"
	"time"

	"sykell-challenge/backend/models"
	"sykell-challenge/backend/repositories"
	"sykell-challenge/backend/services/socket"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

// CrawlManager manages all active crawl jobs
type CrawlManager struct {
	jobs         map[string]*CrawlJob
	mu           sync.RWMutex
	urlRepo      *repositories.URLRepository
	crawlJobRepo *repositories.CrawlJobRepository
}

var crawlManager *CrawlManager

// InitializeCrawlManager initializes the global crawl manager with database connection
func InitializeCrawlManager(db *gorm.DB) {
	crawlManager = &CrawlManager{
		jobs:         make(map[string]*CrawlJob),
		urlRepo:      repositories.NewURLRepository(db),
		crawlJobRepo: repositories.NewCrawlJobRepository(db),
	}
}

// CrawlJobOrExistingData represents either a new crawl job or existing URL data
type CrawlJobOrExistingData struct {
	Job          *CrawlJob   `json:"job,omitempty"`
	ExistingData *models.URL `json:"existing_data,omitempty"`
	IsExisting   bool        `json:"is_existing"`
}

// StartCrawlJob creates and starts a new crawl job or returns existing data
func (cm *CrawlManager) StartCrawlJob(url string) (*CrawlJob, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.urlRepo == nil || cm.crawlJobRepo == nil {
		return nil, fmt.Errorf("crawl manager not initialized with database")
	}

	// Check if URL already exists in database
	existingURL, err := cm.urlRepo.GetByURL(url)
	if err == nil {
		// URL exists, check if it has been successfully crawled
		if existingURL.Status == "done" && existingURL.Title != "" {
			// Return existing data without starting a new crawl
			fmt.Printf("URL %s already exists with completed data, returning existing data\n", url)

			// Create a mock job structure to maintain API compatibility
			// This won't actually run a crawl, but provides the expected response format
			mockJob := &CrawlJob{
				ID:        existingURL.CrawlJobID,
				URL:       url,
				URLID:     existingURL.ID,
				Status:    "completed",
				StartTime: existingURL.UpdatedAt,
				Done:      make(chan bool),
				Cancel:    make(chan bool),
				Result:    make(chan CrawlResult, 1),
			}

			// Send the existing data through the result channel
			go func() {
				defer func() {
					close(mockJob.Done)
					close(mockJob.Cancel)
					close(mockJob.Result)
				}()

				mockJob.Result <- CrawlResult{
					Success: true,
					Data:    *existingURL,
				}
			}()

			return mockJob, nil
		} else if existingURL.Status == "running" || existingURL.Status == "queued" {
			// URL is currently being crawled, return error or the existing job info
			return nil, fmt.Errorf("URL %s is already being crawled (status: %s)", url, existingURL.Status)
		}
		// If status is "error", we'll continue to re-crawl
	}

	jobID := uuid.New().String()

	// Create URL record in database with crawl job ID
	var urlRecord models.URL
	if existingURL != nil {
		// Update existing record for re-crawling
		urlRecord = *existingURL
		urlRecord.Status = "running"
		urlRecord.CrawlJobID = jobID
		if err := cm.urlRepo.Update(&urlRecord); err != nil {
			return nil, fmt.Errorf("failed to update existing URL record: %v", err)
		}
	} else {
		// Create new URL record
		urlRecord = models.URL{
			URL:        url,
			Status:     "running",
			CrawlJobID: jobID,
		}
		if err := cm.urlRepo.Create(&urlRecord); err != nil {
			return nil, fmt.Errorf("failed to create URL record: %v", err)
		}
	}

	// Create crawl job record in database
	now := time.Now()
	crawlJobRecord := models.CrawlJob{
		JobID:     jobID,
		URL:       url,
		URLID:     urlRecord.ID,
		Status:    "queued",
		StartedAt: nil, // Will be set when crawl actually starts
	}

	// Manually set timestamps to avoid zero-value issues
	crawlJobRecord.CreatedAt = now
	crawlJobRecord.UpdatedAt = now

	if err := cm.crawlJobRepo.Create(&crawlJobRecord); err != nil {
		return nil, fmt.Errorf("failed to create crawl job record: %v", err)
	}

	// Broadcast job queued
	socket.BroadcastCrawlUpdate("crawl_queued", map[string]interface{}{
		"jobId":  jobID,
		"url":    url,
		"status": "queued",
		"url_id": urlRecord.ID,
	})

	job := &CrawlJob{
		ID:        jobID,
		URL:       url,
		URLID:     urlRecord.ID,
		Status:    "running",
		StartTime: now,
		Done:      make(chan bool),
		Cancel:    make(chan bool),
		Result:    make(chan CrawlResult, 1),
	}

	cm.jobs[jobID] = job

	// Start the crawl routine
	go func() {
		defer func() {
			// Clean up job from memory
			cm.mu.Lock()
			delete(cm.jobs, jobID)
			cm.mu.Unlock()

			// Close channels safely
			select {
			case job.Done <- true:
			default:
			}
			close(job.Done)
			close(job.Result)

			// Close cancel channel if it's not already closed
			select {
			case <-job.Cancel:
				// Channel was already closed or has a value
			default:
				close(job.Cancel)
			}
		}()

		crawlJobRecord.SetStarted()
		cm.crawlJobRepo.Update(&crawlJobRecord)

		socket.BroadcastCrawlUpdate("crawl_started", map[string]interface{}{
			"jobId":  jobID,
			"url":    url,
			"status": "running",
			"url_id": job.URLID,
		})

		// Create a done channel to signal completion or cancellation
		done := make(chan bool, 1)
		var crawlData CrawlData

		// Start the actual crawl work in a separate goroutine
		go func() {
			crawlData = CrawlURLWithCallback(url)
			done <- true
		}()

		// Wait for either completion or cancellation
		select {
		case <-job.Cancel:
			// Job was cancelled
			job.Status = "cancelled"
			crawlJobRecord.SetCancelled()
			cm.crawlJobRepo.Update(&crawlJobRecord)
			cm.urlRepo.UpdateStatus(job.URLID, "error")

			socket.BroadcastCrawlUpdate("crawl_cancelled", map[string]interface{}{
				"jobId":  jobID,
				"url":    url,
				"status": "cancelled",
				"url_id": job.URLID,
			})

			job.Result <- CrawlResult{Success: false, Error: "Job cancelled"}
			return

		case <-done:
			// Crawl completed, continue with processing
		}

		// Fetch the existing URL record to preserve timestamps
		existingURL, err := cm.urlRepo.GetByID(job.URLID)
		if err != nil {
			job.Status = "error"
			crawlJobRecord.SetError(fmt.Sprintf("Failed to fetch existing URL record: %v", err))
			cm.crawlJobRepo.Update(&crawlJobRecord)

			socket.BroadcastCrawlUpdate("crawl_error", map[string]interface{}{
				"jobId":  jobID,
				"url":    url,
				"status": "error",
				"error":  fmt.Sprintf("Failed to fetch existing URL record: %v", err),
				"url_id": job.URLID,
			})

			job.Result <- CrawlResult{Success: false, Error: fmt.Sprintf("Failed to fetch existing URL record: %v", err)}
			return
		}

		// Update with main collector data (half-completed)
		mainData := crawlData.MainData
		existingURL.Title = mainData.Title
		existingURL.StatusCode = mainData.StatusCode
		existingURL.HTMLVersion = mainData.HTMLVersion
		existingURL.LoginForm = mainData.LoginForm
		existingURL.Tags = mainData.Tags
		// Don't update Links yet - they're not processed

		// Emit half-completed event with main data
		socket.BroadcastCrawlUpdate("crawl_half_completed", map[string]interface{}{
			"jobId":        jobID,
			"url":          url,
			"url_id":       job.URLID,
			"title":        existingURL.Title,
			"status_code":  existingURL.StatusCode,
			"html_version": existingURL.HTMLVersion,
			"login_form":   existingURL.LoginForm,
			"tags_count":   len(existingURL.Tags),
			"links_count":  crawlData.LinkCount,
			"status":       "half_completed",
		})

		// Check for cancellation before processing links
		select {
		case <-job.Cancel:
			job.Status = "cancelled"
			crawlJobRecord.SetCancelled()
			cm.crawlJobRepo.Update(&crawlJobRecord)
			cm.urlRepo.UpdateStatus(job.URLID, "error")

			socket.BroadcastCrawlUpdate("crawl_cancelled", map[string]interface{}{
				"jobId":  jobID,
				"url":    url,
				"status": "cancelled",
				"url_id": job.URLID,
			})

			job.Result <- CrawlResult{Success: false, Error: "Job cancelled"}
			return
		default:
			// Continue processing links if not cancelled
		}

		// Now process the links with secondary collector
		fullData := crawlData.ProcessLinks()

		// Update with full data including processed links
		existingURL.Links = fullData.Links
		existingURL.Status = "done"
		existingURL.CrawlJobID = jobID

		// Update the database record with complete crawl results
		if err := cm.urlRepo.Update(existingURL); err != nil {
			job.Status = "error"
			crawlJobRecord.SetError(fmt.Sprintf("Failed to update database: %v", err))
			cm.crawlJobRepo.Update(&crawlJobRecord)

			socket.BroadcastCrawlUpdate("crawl_error", map[string]interface{}{
				"jobId":  jobID,
				"url":    url,
				"status": "error",
				"error":  fmt.Sprintf("Failed to update database: %v", err),
				"url_id": job.URLID,
			})

			job.Result <- CrawlResult{Success: false, Error: fmt.Sprintf("Failed to update database: %v", err)}
			return
		}

		job.Status = "completed"
		crawlJobRecord.SetCompleted()
		cm.crawlJobRepo.Update(&crawlJobRecord)

		socket.BroadcastCrawlUpdate("crawl_completed", fullData)

		job.Result <- CrawlResult{Success: true, Data: *existingURL}
	}()

	return job, nil
}

// GetJob retrieves a job by ID
func (cm *CrawlManager) GetJob(jobID string) (*CrawlJob, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	job, exists := cm.jobs[jobID]
	return job, exists
}

// CancelJob cancels a running job
func (cm *CrawlManager) CancelJob(jobID string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	job, exists := cm.jobs[jobID]
	if !exists {
		return false
	}

	// Try to send cancellation signal
	select {
	case job.Cancel <- true:
		fmt.Printf("Cancellation signal sent for job %s\n", jobID)
		return true
	default:
		// Channel might be full or closed, but job might still be running
		fmt.Printf("Could not send cancellation signal for job %s (channel full/closed)\n", jobID)
		return false
	}
}

// ListActiveJobs returns all currently active jobs
func (cm *CrawlManager) ListActiveJobs() map[string]*CrawlJob {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	jobs := make(map[string]*CrawlJob)
	for id, job := range cm.jobs {
		jobs[id] = job
	}
	return jobs
}

// GetCrawlJobRepo returns the crawl job repository
func (cm *CrawlManager) GetCrawlJobRepo() *repositories.CrawlJobRepository {
	return cm.crawlJobRepo
}

// GetURLRepo returns the URL repository
func (cm *CrawlManager) GetURLRepo() *repositories.URLRepository {
	return cm.urlRepo
}

// GetCrawlManager returns the global crawl manager (for testing purposes)
func GetCrawlManager() *CrawlManager {
	return crawlManager
}

// GetActiveJobsCount returns the number of currently active crawl jobs
func GetActiveJobsCount() int {
	if crawlManager == nil {
		return 0
	}
	crawlManager.mu.RLock()
	defer crawlManager.mu.RUnlock()
	return len(crawlManager.jobs)
}

// CleanupOldJobs removes completed job records older than the specified duration
func CleanupOldJobs(olderThanDays int) error {
	if crawlManager == nil {
		return fmt.Errorf("crawl manager not initialized")
	}

	cutoffTime := time.Now().AddDate(0, 0, -olderThanDays)
	return crawlManager.crawlJobRepo.DeleteOldJobs(cutoffTime)
}

// GetJobStats returns statistics about crawl jobs
func GetJobStats() (map[string]interface{}, error) {
	if crawlManager == nil {
		return nil, fmt.Errorf("crawl manager not initialized")
	}

	stats := map[string]interface{}{}

	// Get active jobs count
	activeJobs, err := crawlManager.crawlJobRepo.GetActiveJobs()
	if err != nil {
		return nil, fmt.Errorf("failed to get active jobs: %v", err)
	}
	stats["active_jobs"] = len(activeJobs)

	// Get jobs by status
	statuses := []string{"queued", "running", "completed", "cancelled", "error"}
	statusCounts := make(map[string]int)

	for _, status := range statuses {
		jobs, err := crawlManager.crawlJobRepo.GetJobsByStatus(status)
		if err != nil {
			return nil, fmt.Errorf("failed to get jobs for status %s: %v", status, err)
		}
		statusCounts[status] = len(jobs)
	}
	stats["status_counts"] = statusCounts

	// Get jobs created in last 24 hours
	last24h := time.Now().Add(-24 * time.Hour)
	recentJobs, err := crawlManager.crawlJobRepo.GetJobsCreatedAfter(last24h)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent jobs: %v", err)
	}
	stats["jobs_last_24h"] = len(recentJobs)

	// Memory jobs count
	stats["memory_jobs"] = GetActiveJobsCount()

	return stats, nil
}
