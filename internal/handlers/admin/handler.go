package admin

import (
	"context"
	"errors"
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
	err := h.cache.SetState(ctx, strconv.FormatInt(message.Chat.ID, 10), "hello there")
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	return tgbotapi.MessageConfig{}, errors.New("no message")
}
