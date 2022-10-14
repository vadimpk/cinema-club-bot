package cache

import "context"

type Cache interface {
	SetState(ctx context.Context, chatID string, state string) error
	GetState(ctx context.Context, chatID string) (string, error)
	SetAdminState(ctx context.Context, chatID string, state string) error
	GetAdminState(ctx context.Context, chatID string) (string, error)
	SetIdentifier(ctx context.Context, chatID, identifier string) error
	GetIdentifier(ctx context.Context, chatID string) (string, error)
	RemoveIdentifier(ctx context.Context, chatID string) error
	SetName(ctx context.Context, chatID string, name string) error
	GetName(ctx context.Context, chatID string) (string, error)
}
