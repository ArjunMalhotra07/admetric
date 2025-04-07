package model

import (
	"time"

	"gorm.io/gorm"
)

type Ad struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	ImageURL  string         `json:"image_url" gorm:"not null"`
	TargetURL string         `json:"target_url" gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Adding this is optional but helps with querying
	Clicks []Click `gorm:"foreignKey:AdID"`
}
