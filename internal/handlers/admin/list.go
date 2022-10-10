package admin

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"strconv"
)

func (h *Handler) updateList(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(list *domain.List, message *tgbotapi.Message) error,
	nextState string, successTextMessage string) (tgbotapi.MessageConfig, error) {

	// get identifier from cache
	identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when reading cache:", err, message.Chat.ID)
	}
	// get event from db
	event, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		return h.errorDB("Unexpected error when reading event:", err, message.Chat.ID)
	}

	// get list from db
	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when reading event:", err, message.Chat.ID)
	}

	err = updateFunc(&list, message)
	if err != nil {
		return tgbotapi.NewMessage(message.Chat.ID, err.Error()), err
	}

	// update entry in db
	err = h.repos.UpdateList(ctx, list)
	if err != nil {
		return h.errorDB("Unexpected error when updating event:", err, message.Chat.ID)
	}

	// set new state
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), nextState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	return tgbotapi.NewMessage(message.Chat.ID, successTextMessage), nil
}

func (h *Handler) updateListCapacity(list *domain.List, message *tgbotapi.Message) error {
	i, err := strconv.Atoi(message.Text)
	if err != nil {
		return err
	}
	list.Capacity = i
	return nil
}
