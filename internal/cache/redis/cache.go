package redis

import (
	"github.com/go-redis/redis/v9"
	"github.com/vadimpk/cinema-club-bot/internal/config"
)

type Cache struct {
	db *redis.Client
}

func NewCache(cfg config.RedisConfig) *Cache {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.URL + cfg.Port,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	return &Cache{db: rdb}
}
