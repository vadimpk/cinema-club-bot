package admin

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func (h *Handler) createNewEvent(ctx context.Context, message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {

	// check if identifier is unique
	identifier := message.Text
	_, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
		}
	} else {
		return tgbotapi.NewMessage(message.Chat.ID, "Подія з таким ідентифікатором вже існує. Введіть новий ідентифікатор:"), nil
	}

	id, err := h.repos.CreateList(ctx, domain.List{
		EventIdentifier: identifier,
	})
	if err != nil {
		return h.errorDB("Unexpected error when creating event: ", err, message.Chat.ID)
	}

	err = h.repos.CreateEvent(ctx, domain.Event{
		Identifier: identifier,
		ListID:     id,
	})
	if err != nil {
		return h.errorDB("Unexpected error when creating event: ", err, message.Chat.ID)
	}

	// set state to cache
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), updateNameState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// set identifier to cache
	err = h.cache.SetIdentifier(ctx, convertChatIDToString(message.Chat.ID), identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	return tgbotapi.NewMessage(message.Chat.ID, "Введіть назву події: "), nil
}

func (h *Handler) updateEvent(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(event *domain.Event, message *tgbotapi.Message) error,
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

	err = updateFunc(&event, message)
	if err != nil {
		return tgbotapi.NewMessage(message.Chat.ID, err.Error()), err
	}

	// update entry in db
	err = h.repos.UpdateEvent(ctx, event)
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

func (h *Handler) updateEventName(event *domain.Event, message *tgbotapi.Message) error {
	event.Name = message.Text
	return nil
}

func (h *Handler) updateEventDescription(event *domain.Event, message *tgbotapi.Message) error {
	event.Description = message.Text
	return nil
}

func (h *Handler) updateEventDate(event *domain.Event, message *tgbotapi.Message) error {
	date, err := time.Parse(time.RFC3339, message.Text+"Z")
	if err != nil {
		log.Println(err)
		return errors.New("Ви ввели неправильний формат дати, спробуйте ще раз: ")
	}
	event.Date = date
	return nil
}

func (h *Handler) updateEventActiveStatus(event *domain.Event, message *tgbotapi.Message) error {
	if message.Text == "yes" {
		event.Active = true
	}
	return nil
}
