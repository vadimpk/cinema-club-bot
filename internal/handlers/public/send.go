package public

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

func (h *Handler) sendMessagesFromAdmin(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	admin, err := h.repos.GetAdmin(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "У вас немає прав на команду /send")}
		}
		return h.errorDB("Unexpected error getting cache", err, message.Chat.ID)
	}

	messages := admin.Messages
	messagesToSend := make([]tgbotapi.MessageConfig, 0)
	for _, m := range messages {
		chatID, err := strconv.Atoi(m.ChatID)
		if err == nil {
			messagesToSend = append(messagesToSend, tgbotapi.NewMessage(int64(chatID), m.Text))
		}
	}

	err = h.repos.ClearAdminMessages(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		messagesToSend = append(messagesToSend, tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%d на сервері сталася помилка і надіслані повідомлення не було видалено з бази даних, напишіть тому, хто робив бота", len(messagesToSend))))
	}
	messagesToSend = append(messagesToSend, tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%d повідомлень надіслано", len(messagesToSend))))

	return messagesToSend

}
