package domain

type Admin struct {
	ChatID   string    `bson:"chat_id" json:"chat_id"`
	Messages []Message `bson:"messages" json:"messages"`
}
