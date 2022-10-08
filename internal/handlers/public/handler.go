package public

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"strconv"
)

type Handler struct {
	cache cache.Cache
}

func NewHandler(cache cache.Cache) *Handler {
	return &Handler{cache: cache}
}

func (h *Handler) HandleMessage(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	var ctx = context.Background()
	state, err := h.cache.GetState(ctx, strconv.FormatInt(message.Chat.ID, 10))
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, state)
	return msg, nil
}
