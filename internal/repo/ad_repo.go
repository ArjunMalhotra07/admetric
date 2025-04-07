package repo

import (
	"database/sql"
	"log"

	"github.com/ArjunMalhotra/internal/model"
)

type AdRepo struct {
	db *sql.DB
}

func NewAdRepository(db *sql.DB) *AdRepo {
	return &AdRepo{db: db}
}

func (r *AdRepo) FetchAll() ([]model.Ad, error) {
	query := `SELECT * FROM ads`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []model.Ad
	for rows.Next() {
		var ad model.Ad
		if err := rows.Scan(&ad.ID, &ad.ImageURL, &ad.TargetURL); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}
	return ads, nil
}

func (r *AdRepo) CountAds() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM ads`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Failed to count ads: %v", err)
		return 0, err
	}
	return count, nil
}
