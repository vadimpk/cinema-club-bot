package domain

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Capacity        int                `bson:"capacity" json:"capacity"`
	EventIdentifier string             `bson:"event_identifier" json:"event_identifier"`
	List            []User             `bson:"list" json:"list"`
}

func (l List) Preview() string {
	text := fmt.Sprintf("Подія: %s\n\n", l.EventIdentifier)
	for i, u := range l.List {
		text += fmt.Sprintf("%d. %s\n%s\n@%s\n\n", i+1, u.Name, u.Phone, u.Username)
	}
	return text + fmt.Sprintf("Всього зареєстровано: %d / %d", len(l.List), l.Capacity)
}
