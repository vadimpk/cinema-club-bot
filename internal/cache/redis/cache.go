package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"time"
)

type Cache struct {
	db  *redis.Client
	ttl time.Duration
}

func NewCache(cfg config.RedisConfig) *Cache {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.URL + cfg.Port,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	return &Cache{db: rdb, ttl: cfg.TTL}
}

func (c *Cache) SetState(ctx context.Context, chatID string, state string) error {
	return c.db.Set(ctx, chatID, state, c.ttl).Err()
}

func (c *Cache) GetState(ctx context.Context, chatID string) (string, error) {
	return c.db.Get(ctx, chatID).Result()
}
