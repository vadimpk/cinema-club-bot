package admin

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"github.com/vadimpk/cinema-club-bot/internal/repository"
	"log"
	"strconv"
)

type Handler struct {
	cache cache.Cache
	repos repository.Repositories
}

const (
	startState             = "starting state"
	createState            = "creating new event state"
	updateNameState        = "updating name of the event"
	updateDescriptionState = "updating description of the event"
	updateDateState        = "updating state of the event"
	updateActiveState      = "updating active status of the event"
)

func NewHandler(cache cache.Cache, repos repository.Repositories) *Handler {
	return &Handler{cache: cache, repos: repos}
}

func (h *Handler) HandleMessage(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	var ctx = context.Background()
	chatID := convertChatIDToString(message.Chat.ID)

	state, err := h.cache.GetState(ctx, chatID)
	if err != nil {
		// if state not found try to init admin
		if err == redis.Nil {
			err := h.initAdmin(ctx, chatID)
			// if unable to init admin, return message saying forbidden
			if err != nil {
				return tgbotapi.NewMessage(message.Chat.ID, err.Error()), err
			}
			state = startState
		} else {
			log.Println("Unexpected error when reading cache: ", err)
		}
	}

	if message.IsCommand() {
		return h.handleCommands(ctx, message)
	}

	switch state {
	case startState:
		return tgbotapi.NewMessage(message.Chat.ID, "Привіт."), nil
	case createState:
		return h.createNewEvent(ctx, message)
	case updateNameState:
		return h.updateEvent(ctx, message, h.updateEventName, updateDescriptionState, "Введіть опис події:")
	case updateDescriptionState:
		return h.updateEvent(ctx, message, h.updateEventDescription, updateDateState, "Введіть дату події (2022-02-22T15:30:00): ")
	case updateDateState:
		return h.updateEvent(ctx, message, h.updateEventDate, updateActiveState, "Активувати подію одразу (yes):")
	case updateActiveState:
		return h.updateEvent(ctx, message, h.updateEventActiveStatus, startState, "Подію успішно створено.")
	}

	return tgbotapi.MessageConfig{}, err
}

func (h *Handler) initAdmin(ctx context.Context, chatID string) error {
	ok := h.repos.IsAdmin(ctx, chatID)
	if ok {
		return nil
	}
	return errors.New("forbidden")
}

func convertChatIDToString(chatID int64) string {
	return strconv.FormatInt(chatID, 10)
}
