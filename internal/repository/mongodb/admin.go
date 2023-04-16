package mongodb

import (
	"context"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"github.com/vadimpk/cinema-club-bot/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	db     *mongo.Collection
	logger logging.Logger
}

func NewAdminRepository(mdb *mongo.Database, logger logging.Logger) *AdminRepository {
	return &AdminRepository{db: mdb.Collection(adminsCollections), logger: logger}
}

func (r *AdminRepository) IsAdmin(ctx context.Context, chatID string) bool {
	err := r.db.FindOne(ctx, bson.M{"chat_id": chatID}).Err()
	if err != nil {
		r.logger.Error("FAILED TO RETRIEVE VALUES FROM ADMINS COLLECTION ", err)
		return false
	}
	return true
}

func (r *AdminRepository) GetAdmin(ctx context.Context, chatID string) (domain.Admin, error) {
	var admin domain.Admin
	err := r.db.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&admin)
	return admin, err
}
