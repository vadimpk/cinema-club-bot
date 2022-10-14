package repository

import (
	"context"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admins interface {
	IsAdmin(ctx context.Context, chatID string) bool
	GetAdmin(ctx context.Context, chatID string) (domain.Admin, error)
	AddMessagesToAdmin(ctx context.Context, chatID string, messages []domain.Message) error
	ClearAdminMessages(ctx context.Context, chatID string) error
}

type Events interface {
	CreateEvent(ctx context.Context, obj domain.Event) error
	UpdateEvent(ctx context.Context, obj domain.Event) error
	GetEvent(ctx context.Context, identifier string) (domain.Event, error)
	GetAll(ctx context.Context) ([]domain.Event, error)
	GetActive(ctx context.Context) ([]domain.Event, error)
	DeleteEvent(ctx context.Context, identifier string) error
	CreateTestEvent(ctx context.Context, identifier string, listID primitive.ObjectID) error
}

type Lists interface {
	CreateList(ctx context.Context, obj domain.List) (primitive.ObjectID, error)
	UpdateList(ctx context.Context, obj domain.List) error
	GetList(ctx context.Context, id primitive.ObjectID) (domain.List, error)
	CreateTestList(ctx context.Context, identifier string, chatID string) (primitive.ObjectID, error)
}

type Repositories interface {
	Admins
	Events
	Lists
}
