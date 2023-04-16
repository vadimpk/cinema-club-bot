package public

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (h *Handler) errorDB(msg string, err error, chatID int64) []tgbotapi.MessageConfig {
	logger := h.logger.Named("errorDB")
	logger = logger.With("chatID", chatID)
	logger = logger.With("message", msg)
	logger = logger.With("error", err)
	logger.Error("failed to handle message")
	return []tgbotapi.MessageConfig{tgbotapi.NewMessage(chatID, "Сталася помилка. Натисніть команду /start\nЯкщо проблема не зникає - пишіть @vadimpk")}
}
