package mongodb

import (
	"context"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventsRepository struct {
	db *mongo.Collection
}

func NewEventsRepository(mdb *mongo.Database) *EventsRepository {
	return &EventsRepository{db: mdb.Collection(eventsCollection)}
}

func (r *EventsRepository) CreateEvent(ctx context.Context, obj domain.Event) error {
	_, err := r.db.InsertOne(ctx, obj)
	return err
}

func (r *EventsRepository) UpdateEvent(ctx context.Context, obj domain.Event) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": obj.ID}, bson.M{"$set": obj})
	return err
}

func (r *EventsRepository) GetEvent(ctx context.Context, identifier string) (domain.Event, error) {
	var event domain.Event
	err := r.db.FindOne(ctx, bson.M{"identifier": identifier}).Decode(&event)
	return event, err
}

func (r *EventsRepository) GetAll(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	cur, err := r.db.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &events)

	return events, err
}

func (r *EventsRepository) GetActive(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	cur, err := r.db.Find(ctx, bson.M{"active": true})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &events)

	return events, err
}

func (r *EventsRepository) DeleteEvent(ctx context.Context, identifier string) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"identifier": identifier})
	return err
}
