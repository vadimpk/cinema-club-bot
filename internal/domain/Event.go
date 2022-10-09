package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Identifier  string             `bson:"identifier,omitempty" json:"identifier"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	ListID      primitive.ObjectID `bson:"list_id" json:"list_id"`
	Date        time.Time          `bson:"date" json:"date"`
	Active      bool               `bson:"active" json:"active"`
}
