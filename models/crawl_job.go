package models

import (
	"time"

	"gorm.io/gorm"
)

type CrawlJob struct {
	gorm.Model
	JobID       string     `json:"jobId" gorm:"type:varchar(255);uniqueIndex;not null"`
	URL         string     `json:"url" gorm:"not null"`
	URLID       uint       `json:"url_id" gorm:"index"`
	Status      string     `json:"status" gorm:"type:enum('queued','running','completed','cancelled','error');default:'queued';not null"`
	StartedAt   *time.Time `json:"started_at" gorm:"default:null"`
	CompletedAt *time.Time `json:"completed_at" gorm:"default:null"`
	Duration    *int64     `json:"duration_ms" gorm:"default:null"` // Duration in milliseconds
	ErrorMsg    string     `json:"error_message,omitempty"`
}

// GetRunningDuration returns the duration the job has been running in milliseconds
func (cj *CrawlJob) GetRunningDuration() int64 {
	if cj.StartedAt == nil {
		return 0
	}

	endTime := time.Now()
	if cj.CompletedAt != nil {
		endTime = *cj.CompletedAt
	}

	return endTime.Sub(*cj.StartedAt).Milliseconds()
}

// IsActive returns true if the job is currently running
func (cj *CrawlJob) IsActive() bool {
	return cj.Status == "running" || cj.Status == "queued"
}

// SetCompleted marks the job as completed and sets the completion time
func (cj *CrawlJob) SetCompleted() {
	now := time.Now()
	cj.CompletedAt = &now
	cj.Status = "completed"
	if cj.StartedAt != nil {
		duration := now.Sub(*cj.StartedAt).Milliseconds()
		cj.Duration = &duration
	}
}

// SetError marks the job as errored and sets the completion time
func (cj *CrawlJob) SetError(errorMsg string) {
	now := time.Now()
	cj.CompletedAt = &now
	cj.Status = "error"
	cj.ErrorMsg = errorMsg
	if cj.StartedAt != nil {
		duration := now.Sub(*cj.StartedAt).Milliseconds()
		cj.Duration = &duration
	}
}

// SetCancelled marks the job as cancelled and sets the completion time
func (cj *CrawlJob) SetCancelled() {
	now := time.Now()
	cj.CompletedAt = &now
	cj.Status = "cancelled"
	if cj.StartedAt != nil {
		duration := now.Sub(*cj.StartedAt).Milliseconds()
		cj.Duration = &duration
	}
}

// SetStarted marks the job as started and sets the start time
func (cj *CrawlJob) SetStarted() {
	now := time.Now()
	cj.StartedAt = &now
	cj.Status = "running"
}
