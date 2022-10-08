package cache

import "context"

type Cache interface {
	SetState(ctx context.Context, chatID string, state string) error
	GetState(ctx context.Context, chatID string) (string, error)
}
