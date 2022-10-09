package repository

import (
	"context"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admins interface {
	IsAdmin(ctx context.Context, chatID string) bool
}

type Events interface {
	Create(ctx context.Context, obj domain.Event) error
	Update(ctx context.Context, obj domain.Event) error
	Get(ctx context.Context, identifier string) (domain.Event, error)
	GetAll(ctx context.Context) ([]domain.Event, error)
	GetActive(ctx context.Context) ([]domain.Event, error)
	Delete(ctx context.Context, identifier string) error
}

type Lists interface {
	Create(ctx context.Context, obj domain.List) error
	InsertUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, userPhone string) error
	Get(ctx context.Context, id primitive.ObjectID) ([]domain.List, error)
	GetActive(ctx context.Context) ([]domain.Event, error)
	Delete(ctx context.Context, identifier string) error
}

type Repositories interface {
	Admins
	Events
}
