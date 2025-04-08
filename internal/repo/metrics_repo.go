package repo

import (
	"log"

	"github.com/go-redis/redis"
)

type MetricsRepo struct {
	RedisClient *redis.Client
}

func NewMetricsRepo(redis *redis.Client) *MetricsRepo {
	return &MetricsRepo{RedisClient: redis}
}

func (r *MetricsRepo) IncrementClickCount(adID string) error {
	key := "clicks:" + adID
	if err := r.RedisClient.Incr(key).Err(); err != nil {
		log.Printf("Failed to increment click count: %v", err)
		return err
	}
	return nil
}

func (r *MetricsRepo) GetClickCount(adID string) (int64, error) {
	key := "clicks:" + adID
	count, err := r.RedisClient.Get(key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Printf("Failed to get click count: %v", err)
		return 0, err
	}
	return count, nil
}
