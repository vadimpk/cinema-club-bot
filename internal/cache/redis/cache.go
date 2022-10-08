package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
)

type Cache struct {
	db  *redis.Client
	ttl time.Duration
}

func NewCache(db *redis.Client, ttl time.Duration) *Cache {
	return &Cache{db: db, ttl: ttl}
}

func (c *Cache) SetState(ctx context.Context, chatID string, state string) error {
	return c.db.Set(ctx, chatID, state, c.ttl).Err()
}

func (c *Cache) GetState(ctx context.Context, chatID string) (string, error) {
	return c.db.Get(ctx, chatID).Result()
}
