package cache

import (
	"encoding/json"
	"time"

	"github.com/egnptr/rest-api/entity"
	"github.com/go-redis/redis/v7"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, expires time.Duration) PostCache {
	return &RedisCache{
		host:    host,
		db:      db,
		expires: expires,
	}
}

func (cache *RedisCache) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) Set(key string, value *entity.Post) {
	client := cache.GetClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	client.Set(key, json, cache.expires*time.Second)
}

func (cache *RedisCache) Get(key string) *entity.Post {
	client := cache.GetClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	post := entity.Post{}
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}
	return &post
}
