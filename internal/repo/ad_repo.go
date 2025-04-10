package repo

import (
	"log"

	"github.com/ArjunMalhotra/internal/model"
	"gorm.io/gorm"
)

type AdRepo struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) *AdRepo {
	return &AdRepo{db: db}
}

func (r *AdRepo) FetchAll() ([]model.Ad, error) {
	var ads []model.Ad
	if err := r.db.Preload("Clicks").Find(&ads).Error; err != nil {
		return nil, err
	}
	return ads, nil
}

func (r *AdRepo) CountAds() (int, error) {
	var count int64
	if err := r.db.Model(&model.Ad{}).Count(&count).Error; err != nil {
		log.Printf("Failed to count ads: %v", err)
		return 0, err
	}
	return int(count), nil
}