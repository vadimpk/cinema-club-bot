package mongodb

import (
	"context"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ListsRepository struct {
	db *mongo.Collection
}

func NewListsRepository(mdb *mongo.Database) *ListsRepository {
	return &ListsRepository{db: mdb.Collection(listsCollection)}
}

func (r *ListsRepository) CreateList(ctx context.Context, obj domain.List) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, obj)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), err
}

func (r *ListsRepository) UpdateList(ctx context.Context, obj domain.List) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": obj.ID}, bson.M{"$set": obj})
	return err
}

func (r *ListsRepository) GetList(ctx context.Context, id primitive.ObjectID) (domain.List, error) {
	var list domain.List
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&list)
	return list, err
}
