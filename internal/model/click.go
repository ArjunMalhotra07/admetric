package model

import "time"

type Click struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	AdID         string    `json:"ad_id" gorm:"not null"`
	Ad           Ad        `gorm:"foreignKey:AdID;references:ID"` // Proper foreign key
	IP           string    `json:"ip" gorm:"not null"`
	PlaybackTime int       `json:"playback_time" gorm:"not null"`
	Timestamp    time.Time `json:"timestamp" gorm:"not null"`
}
