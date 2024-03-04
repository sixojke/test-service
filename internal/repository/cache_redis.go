package repository

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/sixojke/test-service/internal/config"
	"github.com/sixojke/test-service/internal/domain"
)

type CahcheRedis struct {
	db     *redis.Client
	config config.CacheConfig
}

func NewCacheRedis(db *redis.Client, config config.CacheConfig) *CahcheRedis {
	return &CahcheRedis{
		db:     db,
		config: config,
	}
}

func (r *CahcheRedis) SetList(limit, offset int, list *domain.List) error {
	listJSON, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("json marshal cache: %v", err)
	}

	key := generateKey(limit, offset)
	logrus.Info(key)
	if err := r.db.Set(key, listJSON, r.config.Expiration).Err(); err != nil {
		return fmt.Errorf("add cache: %v", err)
	}

	for _, item := range list.Goods {
		if err := r.db.Set(strconv.Itoa(item.Id), key, r.config.Expiration).Err(); err != nil {
			return fmt.Errorf("add cache: %v", err)
		}
	}

	return nil
}

func (r *CahcheRedis) GetList(limit, offset int) (*domain.List, error) {
	key := generateKey(limit, offset)
	listJSON, err := r.db.Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, fmt.Errorf("get list cache: %v", err)
		} else {
			return nil, nil
		}
	}

	var list domain.List
	if err := json.Unmarshal([]byte(listJSON), &list); err != nil {
		return nil, fmt.Errorf("unmarshal listJSON: %v", err)
	}

	return &list, nil
}

func (r *CahcheRedis) Delete(itemId int) error {
	key, err := r.db.Get(strconv.Itoa(itemId)).Result()
	if err != nil {
		if err != redis.Nil {
			return fmt.Errorf("get key cache: %v", err)
		} else {
			return nil
		}
	}

	if err := r.db.Del(key).Err(); err != nil {
		return fmt.Errorf("cache del: %v", err)
	}

	return nil
}

func generateKey(limit, offset int) string {
	return fmt.Sprintf("list:%v:%v", limit, offset)
}
