package domain

type Message struct {
	ChatID string `bson:"chat_id" json:"chat_id"`
	Text   string `bson:"text" json:"text"`
}
