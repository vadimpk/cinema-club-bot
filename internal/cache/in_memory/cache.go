package in_memory

import (
	"context"
	"github.com/patrickmn/go-cache"
	cacheInterface "github.com/vadimpk/cinema-club-bot/internal/cache"
	"time"
)

type LocalCache struct {
	cache    cache.Cache
	ttl      time.Duration
	adminttl time.Duration
}

func NewCache(ttl, adminttl time.Duration) cacheInterface.Cache {
	return &LocalCache{
		cache:    *cache.New(ttl, ttl),
		ttl:      ttl,
		adminttl: adminttl,
	}
}

func (l LocalCache) SetState(ctx context.Context, chatID string, state string) error {
	l.cache.Set(chatID, state, l.ttl)
	return nil
}

func (l LocalCache) GetState(ctx context.Context, chatID string) (string, error) {
	state, ok := l.cache.Get(chatID)
	if !ok {
		return "", nil
	}
	return state.(string), nil
}

func (l LocalCache) SetAdminState(ctx context.Context, chatID string, state string) error {
	l.cache.Set(chatID+"_admin", state, l.adminttl)
	return nil
}

func (l LocalCache) GetAdminState(ctx context.Context, chatID string) (string, error) {
	state, ok := l.cache.Get(chatID + "_admin")
	if !ok {
		return "", nil
	}
	return state.(string), nil
}

func (l LocalCache) SetIdentifier(ctx context.Context, chatID, identifier string) error {
	l.cache.Set(chatID+"_identifier", identifier, l.ttl)
	return nil
}

func (l LocalCache) GetIdentifier(ctx context.Context, chatID string) (string, error) {
	state, ok := l.cache.Get(chatID + "_identifier")
	if !ok {
		return "", nil
	}
	return state.(string), nil
}

func (l LocalCache) SetAdminIdentifier(ctx context.Context, chatID, identifier string) error {
	l.cache.Set(chatID+"_identifier", identifier, l.adminttl)
	return nil
}

func (l LocalCache) RemoveIdentifier(ctx context.Context, chatID string) error {
	l.cache.Delete(chatID + "_identifier")
	return nil
}

func (l LocalCache) SetName(ctx context.Context, chatID string, name string) error {
	l.cache.Set(chatID+"_name", name, l.ttl)
	return nil
}

func (l LocalCache) GetName(ctx context.Context, chatID string) (string, error) {
	state, ok := l.cache.Get(chatID + "_name")
	if !ok {
		return "", nil
	}
	return state.(string), nil
}
