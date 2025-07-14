package models

import (
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	URL         string `json:"url" gorm:"not null"`
	Title       string `json:"title" gorm:"type:varchar(500)"` // Page title
	Status      string `json:"status" gorm:"type:enum('queued','running','done','error', 'cancelled');default:'queued';not null"`
	StatusCode  int    `json:"statusCode" gorm:"default:0"` // HTTP status code (200, 404, 500, etc.)
	HTMLVersion string `json:"htmlVersion"`
	LoginForm   bool   `json:"loginFormPresent" gorm:"default:false"`
	Tags        Tags   `json:"tags" gorm:"type:json"`
	Links       Links  `json:"links" gorm:"type:json"`
	JobId       string `json:"jobId" gorm:"index"` // ID of the channel/goroutine running the crawl
}
