package repo

import (
	"log"
	"time"

	"github.com/ArjunMalhotra/internal/model"
	"gorm.io/gorm"
)

type ClickRepo struct {
	DB *gorm.DB
}

func NewClickRepo(db *gorm.DB) *ClickRepo {
	return &ClickRepo{DB: db}
}

func (r *ClickRepo) Save(click model.Click) error {
	if err := r.DB.Create(&click).Error; err != nil {
		log.Printf("Failed to save click event: %v", err)
		return err
	}
	return nil
}

func (r *ClickRepo) AdExists(adID string) (bool, error) {
	var exists bool
	err := r.DB.Model(&model.Ad{}).
		Select("count(?) > 0").
		Where("id = ?", adID).
		Scan(&exists).Error

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ClickRepo) GetClickCountByIP(ip string) (int, error) {
	var count int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	err := r.DB.Model(&model.Click{}).Where("ip = ? AND timestamp > ?", ip, oneHourAgo).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
