package domain

type User struct {
	Name   string `bson:"name" json:"name"`
	Phone  string `bson:"phone" json:"phone"`
	ChatID string `bson:"chat_id" json:"chat_id"`
}
