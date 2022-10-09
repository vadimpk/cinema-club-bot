package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type AdminRepository struct {
	db *mongo.Collection
}

func NewAdminRepository(mdb *mongo.Database) *AdminRepository {
	return &AdminRepository{db: mdb.Collection(adminsCollections)}
}

func (r *AdminRepository) IsAdmin(ctx context.Context, chatID string) bool {
	err := r.db.FindOne(ctx, bson.M{"chat_id": chatID}).Err()
	if err != nil {
		log.Println("FAILED TO RETRIEVE VALUES FROM ADMINS COLLECTION ", err)
		return false
	}
	return true
}
