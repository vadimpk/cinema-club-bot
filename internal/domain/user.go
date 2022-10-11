package domain

type User struct {
	Name     string `bson:"name" json:"name"`
	Phone    string `bson:"phone" json:"phone"`
	ChatID   string `bson:"chat_id" json:"chat_id"`
	UserID   string `bson:"user_id" json:"user_id"`
	Username string `bson:"username" json:"username"`
}
