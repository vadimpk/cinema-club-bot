package mongodb

import (
	"context"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
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

func (r *AdminRepository) GetAdmin(ctx context.Context, chatID string) (domain.Admin, error) {
	var admin domain.Admin
	err := r.db.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&admin)
	return admin, err
}

func (r *AdminRepository) AddMessagesToAdmin(ctx context.Context, chatID string, messages []domain.Message) error {
	admin, err := r.GetAdmin(ctx, chatID)
	if err != nil {
		return err
	}
	if admin.Messages == nil {
		admin.Messages = messages
	} else {
		admin.Messages = append(admin.Messages, messages...)
	}
	_, err = r.db.UpdateOne(ctx, bson.M{"chat_id": chatID}, bson.M{"$set": admin})
	return err
}

func (r *AdminRepository) ClearAdminMessages(ctx context.Context, chatID string) error {
	admin, err := r.GetAdmin(ctx, chatID)
	if err != nil {
		return err
	}
	admin.Messages = nil
	_, err = r.db.UpdateOne(ctx, bson.M{"chat_id": chatID}, bson.M{"$set": admin})
	return err
}
