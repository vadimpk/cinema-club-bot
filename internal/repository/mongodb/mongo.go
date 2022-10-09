package mongodb

import (
	"context"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const timeout = 10 * time.Second

type Repositories struct {
	*AdminRepository
	*EventsRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		AdminRepository:  NewAdminRepository(db),
		EventsRepository: NewEventsRepository(db),
	}
}

// NewClient established connection to a mongoDb instance using provided URI and auth credentials.
func NewClient(cfg config.MongoConfig) (*mongo.Client, error) {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
