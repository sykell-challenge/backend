package models

import (
	"time"

	"gorm.io/gorm"
)

type CrawlJob struct {
	gorm.Model
	URL         string     `json:"url" gorm:"not null"`
	URLID       uint       `json:"urlId" gorm:"index"`
	Status      string     `json:"status" gorm:"type:enum('queued','running','completed','cancelled','error');default:'queued';not null"`
	StartedAt   *time.Time `json:"startedAt" gorm:"default:null"`
	CompletedAt *time.Time `json:"completedAt" gorm:"default:null"`
	Progress    int        `json:"progress" gorm:"default:0"` // Progress percentage
	ErrorMsg    string     `json:"errorMessage,omitempty"`
}
