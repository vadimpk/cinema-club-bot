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

func (c *Cache) SetAdminState(ctx context.Context, chatID string, state string) error {
	return c.db.Set(ctx, chatID+"_admin", state, c.ttl).Err()
}

func (c *Cache) GetAdminState(ctx context.Context, chatID string) (string, error) {
	return c.db.Get(ctx, chatID+"_admin").Result()
}

func (c *Cache) SetIdentifier(ctx context.Context, chatID, identifier string) error {
	return c.db.Set(ctx, chatID+"_identifier", identifier, c.ttl).Err()
}

func (c *Cache) GetIdentifier(ctx context.Context, chatID string) (string, error) {
	return c.db.Get(ctx, chatID+"_identifier").Result()
}

func (c *Cache) RemoveIdentifier(ctx context.Context, chatID string) error {
	return c.db.Del(ctx, chatID+"_identifier").Err()
}

func (c *Cache) SetName(ctx context.Context, chatID string, name string) error {
	return c.db.Set(ctx, chatID+"_name", name, c.ttl).Err()
}

func (c *Cache) GetName(ctx context.Context, chatID string) (string, error) {
	return c.db.Get(ctx, chatID+"_name").Result()
}
