package admin

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) handleCreationStartProcess(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	err := h.cache.SetAdminState(ctx, convertChatIDToString(message.Chat.ID), createState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введіть ідентифікатор нової події:")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return []tgbotapi.MessageConfig{msg}
}

func (h *Handler) createNewEvent(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {

	// check if identifier is unique
	identifier := message.Text
	_, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
		}
	} else {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Подія з таким ідентифікатором вже існує. Введіть новий ідентифікатор:")}
	}

	// create list
	id, err := h.repos.CreateList(ctx, domain.List{
		EventIdentifier: identifier,
	})
	if err != nil {
		return h.errorDB("Unexpected error when creating event: ", err, message.Chat.ID)
	}

	// create event
	err = h.repos.CreateEvent(ctx, domain.Event{
		Identifier: identifier,
		ListID:     id,
	})
	if err != nil {
		return h.errorDB("Unexpected error when creating event: ", err, message.Chat.ID)
	}

	// set state to cache
	err = h.cache.SetAdminState(ctx, convertChatIDToString(message.Chat.ID), updateNameStateOnCreation)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// set identifier to cache
	err = h.cache.SetIdentifier(ctx, convertChatIDToString(message.Chat.ID), identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Введіть назву події: ")}
}

/*
updateEvent() receives function that updates some field in event (updateFunc) and
nextFunc that gets called after the event is updated. So this function gets identifier
from cache, gets event from repos, updates event using given updateFunc function and
then calls nextFunc
*/
func (h *Handler) updateEvent(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(event *domain.Event, text string) error, replyText string) []tgbotapi.MessageConfig {

	if message.Text == cancelUpdateOption {
		return h.goToMainMenu(ctx, message, "Зміну скасовано.")
	}

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

	err = updateFunc(&event, message.Text)
	if err != nil {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, err.Error())}
	}

	// update entry in db
	err = h.repos.UpdateEvent(ctx, event)
	if err != nil {
		return h.errorDB("Unexpected error when updating event:", err, message.Chat.ID)
	}

	return h.goBackToUpdateOptions(ctx, message, replyText)
}

func (h *Handler) updateEventOnCreation(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(event *domain.Event, text string) error,
	state, replyText string, last bool) []tgbotapi.MessageConfig {

	if message.Text == cancelUpdateOption {
		// get identifier from cache
		identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
		if err != nil {
			return h.errorDB("Unexpected error when reading cache:", err, message.Chat.ID)
		}
		err = h.repos.DeleteEvent(ctx, identifier)
		if err != nil {
			return h.errorDB("Unexpected error when deleting error", err, message.Chat.ID)
		}
		return h.goToMainMenu(ctx, message, "Подію не збережено.")
	}

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

	err = updateFunc(&event, message.Text)
	if err != nil {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, err.Error())}
	}

	// update entry in db
	err = h.repos.UpdateEvent(ctx, event)
	if err != nil {
		return h.errorDB("Unexpected error when updating event:", err, message.Chat.ID)
	}
	if last {
		return h.goToMainMenu(ctx, message, replyText)
	}
	return h.askToEnterData(ctx, message, state, replyText)

}

// helping function to update name of the event. Called from updateEvent
func (h *Handler) updateEventName(event *domain.Event, text string) error {
	event.Name = text
	return nil
}

// helping function to update description of the event. Called from updateEvent
func (h *Handler) updateEventDescription(event *domain.Event, text string) error {
	event.Description = text
	return nil
}

// helping function to update date of the event. Called from updateEvent
func (h *Handler) updateEventDate(event *domain.Event, text string) error {
	invalidError := errors.New("Ви ввели неправильний формат дати, спробуйте ще раз: ")
	parts := strings.Split(text, " ")
	if len(parts) != 5 {
		return invalidError
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return invalidError
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return invalidError
	}
	if m < 1 || m > 12 {
		return invalidError
	}
	var month = time.Month(m)
	hour, err := strconv.Atoi(parts[2])
	if err != nil {
		return invalidError
	}
	minute, err := strconv.Atoi(parts[3])
	if err != nil {
		return invalidError
	}
	timeZone, err := strconv.Atoi(parts[4])
	if err != nil {
		return invalidError
	}
	if timeZone < -12 || timeZone > 12 {
		return invalidError
	}
	location := time.FixedZone("UTC"+parts[4], 0)

	date := time.Date(time.Now().Year(), month, day, hour, minute, 0, 0, location)

	event.Date = date
	return nil
}

// helping function to update active status of the event. Called from updateEvent
func (h *Handler) updateEventActiveStatus(event *domain.Event, text string) error {
	if text == activateEventOption {
		event.Active = true
	} else if text == deactivateEventOption {
		event.Active = false
	}
	return nil
}

func (h *Handler) deleteEvent(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	// get identifier from cache
	identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when reading cache:", err, message.Chat.ID)
	}

	err = h.repos.DeleteEvent(ctx, identifier)
	if err != nil {
		return h.errorDB("Unexpected error when deleting event:", err, message.Chat.ID)
	}

	// set new state
	err = h.cache.SetAdminState(ctx, convertChatIDToString(message.Chat.ID), startState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// remove identifier
	err = h.cache.RemoveIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Подію видалено")
	msg.ReplyMarkup = h.getOptionsKeyboard(false)
	return []tgbotapi.MessageConfig{msg}
}
