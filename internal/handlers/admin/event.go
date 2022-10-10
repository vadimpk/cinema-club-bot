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

func (h *Handler) handleCreationStartProcess(ctx context.Context, message *tgbotapi.Message) tgbotapi.MessageConfig {
	err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), createState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введіть ідентифікатор нової події:")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return msg
}

func (h *Handler) createNewEvent(ctx context.Context, message *tgbotapi.Message) tgbotapi.MessageConfig {

	// check if identifier is unique
	identifier := message.Text
	_, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
		}
	} else {
		return tgbotapi.NewMessage(message.Chat.ID, "Подія з таким ідентифікатором вже існує. Введіть новий ідентифікатор:")
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
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), updateNameStateOnCreation)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// set identifier to cache
	err = h.cache.SetIdentifier(ctx, convertChatIDToString(message.Chat.ID), identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	return tgbotapi.NewMessage(message.Chat.ID, "Введіть назву події: ")
}

/*
updateEvent() receives function that updates some field in event (updateFunc) and
nextFunc that gets called after the event is updated. So this function gets identifier
from cache, gets event from repos, updates event using given updateFunc function and
then calls nextFunc
*/
func (h *Handler) updateEvent(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(event *domain.Event, message *tgbotapi.Message) error, replyText string) tgbotapi.MessageConfig {

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
		return tgbotapi.NewMessage(message.Chat.ID, err.Error())
	}

	// update entry in db
	err = h.repos.UpdateEvent(ctx, event)
	if err != nil {
		return h.errorDB("Unexpected error when updating event:", err, message.Chat.ID)
	}

	return h.goBackToUpdateOptions(ctx, message, replyText)
}

func (h *Handler) updateEventOnCreation(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(event *domain.Event, message *tgbotapi.Message) error,
	state, replyText string, last bool) tgbotapi.MessageConfig {

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
		return tgbotapi.NewMessage(message.Chat.ID, err.Error())
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
func (h *Handler) updateEventName(event *domain.Event, message *tgbotapi.Message) error {
	event.Name = message.Text
	return nil
}

// helping function to update description of the event. Called from updateEvent
func (h *Handler) updateEventDescription(event *domain.Event, message *tgbotapi.Message) error {
	event.Description = message.Text
	return nil
}

// helping function to update date of the event. Called from updateEvent
func (h *Handler) updateEventDate(event *domain.Event, message *tgbotapi.Message) error {
	date, err := time.Parse(time.RFC3339, message.Text+"Z")
	if err != nil {
		log.Println(err)
		return errors.New("Ви ввели неправильний формат дати, спробуйте ще раз: ")
	}
	event.Date = date
	return nil
}

// helping function to update active status of the event. Called from updateEvent
func (h *Handler) updateEventActiveStatus(event *domain.Event, message *tgbotapi.Message) error {
	if message.Text == activateEventOption {
		event.Active = true
	} else if message.Text == deactivateEventOption {
		event.Active = false
	}
	return nil
}

func (h *Handler) deleteEvent(ctx context.Context, message *tgbotapi.Message) tgbotapi.MessageConfig {
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
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), startState)
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

	return msg
}
