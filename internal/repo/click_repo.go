package repo

import (
	"database/sql"
	"log"
	"net"

	"github.com/ArjunMalhotra/internal/model"
)

type ClickRepo struct {
	db *sql.DB
}

func (r *ClickRepo) Save(click model.Click) error {
	query := `INSERT INTO clicks (ad_id, timestamp, ip, playback_time) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, click.AdID, click.Timestamp, click.IP, click.PlaybackTime)
	if err != nil {
		log.Printf("Failed to save click event: %v", err)
		return err
	}
	return nil
}

func (r *ClickRepo) AdExists(adID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM ads WHERE id = ?)"
	err := r.db.QueryRow(query, adID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ClickRepo) GetClickCountByIP(ip string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM clicks WHERE ip = ? AND timestamp > NOW() - INTERVAL 1 hour"
	err := r.db.QueryRow(query, ip).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ClickRepo) IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func (r *ClickRepo) IsPlaybackTimeValid(playbackTime int) bool {
	return playbackTime > 0
}
