package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type List struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Capacity        int                `bson:"capacity" json:"capacity"`
	EventIdentifier string             `bson:"event_identifier" json:"event_identifier"`
	List            []User             `bson:"list" json:"list"`
}
