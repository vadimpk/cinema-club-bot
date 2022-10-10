package public

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"github.com/vadimpk/cinema-club-bot/internal/repository"
	"strconv"
)

type Handler struct {
	cache cache.Cache
	repos repository.Repositories
}

func NewHandler(cache cache.Cache, repos repository.Repositories) *Handler {
	return &Handler{cache: cache, repos: repos}
}

func (h *Handler) HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var ctx = context.Background()
	state, err := h.cache.GetState(ctx, strconv.FormatInt(message.Chat.ID, 10))
	if err != nil {
		return tgbotapi.MessageConfig{}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, state)
	return msg
}
