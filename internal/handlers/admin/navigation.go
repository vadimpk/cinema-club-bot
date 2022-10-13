package admin

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
chooseUpdateOptions - validates given message from previous step and if valid sets it to cache.
Returns message with keyboard to choose updating options from
*/
func (h *Handler) chooseUpdateOptions(ctx context.Context, message *tgbotapi.Message) tgbotapi.MessageConfig {
	// check if identifier is valid
	identifier := message.Text
	if identifier == toMainMenuOption {
		return h.goToMainMenu(ctx, message, "Виберіть дію:")
	}
	event, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Події з таким ідентифікатором не існує. Виберіть ідентифікатор ще раз:")
	}

	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	// set state to cache
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), chooseEventUpdateOptionsState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// set identifier to cache
	err = h.cache.SetIdentifier(ctx, convertChatIDToString(message.Chat.ID), identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%s\n\n%s", event.Format(list), "Виберіть дію:"))
	msg.ReplyMarkup = h.getUpdateEventOptionsKeyboard(event, true)
	return msg
}

/*
chooseEvent - returns message with keyboard of events' identifiers to choose from
*/
func (h *Handler) chooseEvent(ctx context.Context, message *tgbotapi.Message, nextState string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Оберіть подію: ")
	keyboard, err := h.getEventsKeyboard(ctx, true)
	if err != nil {
		return h.errorDB("Unexpected error when reading events:", err, message.Chat.ID)
	}
	msg.ReplyMarkup = keyboard

	// set state to cache
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), nextState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	return msg
}

/*
askToEnterData - gets called when change of some data is requested (either on creation or via buttons)
Changes state and returns message
*/
func (h *Handler) askToEnterData(ctx context.Context, message *tgbotapi.Message, state, replyText string) tgbotapi.MessageConfig {
	err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), state)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}
	return tgbotapi.NewMessage(message.Chat.ID, replyText)
}

/*
goBackToUpdateOptions - gets called after successful updating of data
to get back to the keyboard with updating options
*/
func (h *Handler) goBackToUpdateOptions(ctx context.Context, message *tgbotapi.Message, replyText string) tgbotapi.MessageConfig {
	err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), chooseEventUpdateOptionsState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
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

	msg := tgbotapi.NewMessage(message.Chat.ID, replyText)
	msg.ReplyMarkup = h.getUpdateEventOptionsKeyboard(event, true)
	return msg
}

/*
goToMainMenu - gets called when user hits toMainMenuButton
*/
func (h *Handler) goToMainMenu(ctx context.Context, message *tgbotapi.Message, msgText string) tgbotapi.MessageConfig {
	// set state to cache
	err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), startState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// remove identifier
	err = h.cache.RemoveIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}
	return h.handleStart(ctx, message, msgText)
}
