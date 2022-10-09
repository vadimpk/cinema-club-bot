package admin

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	createNewCommand = "create_event"
)

func (h *Handler) handleCommands(ctx context.Context, message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	switch message.Command() {
	case createNewCommand:
		err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), createState)
		if err != nil {
			return tgbotapi.MessageConfig{}, err
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Введіть ідентифікатор нової події:"), err
	}
	return tgbotapi.MessageConfig{}, nil
}
