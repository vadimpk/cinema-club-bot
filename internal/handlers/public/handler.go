package public

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMessage(message *tgbotapi.Message) error {
	return nil
}
