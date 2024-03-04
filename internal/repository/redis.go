package repository

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sixojke/test-service/internal/config"
)

func NewRedisDB(cfg config.RedisConfig) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DBName,
	})

	if _, err := r.Ping().Result(); err != nil {
		return nil, fmt.Errorf("connection: %v", err)
	}

	return r, nil
}
