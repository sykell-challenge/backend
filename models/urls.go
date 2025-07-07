package models

import (
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	URL         string `json:"url" gorm:"not null"`
	Status      string `json:"status" gorm:"type:enum('queued','running','done','error');default:'queued';not null"`
	HTMLVersion string `json:"html_version"`
	LoginForm   bool   `json:"login_form" gorm:"default:false"`
	Tags        Tags   `json:"tags" gorm:"type:json"`
	Links       Links  `json:"links" gorm:"type:json"`
}
