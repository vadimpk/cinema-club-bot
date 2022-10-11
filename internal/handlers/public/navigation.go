package public

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

/*
chooseEvent - returns message with keyboard of events' identifiers to choose from
*/
func (h *Handler) chooseEvent(ctx context.Context, message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Оберіть подію: ")
	keyboard, err := h.getEventsKeyboard(ctx, true)
	if err != nil {
		return h.errorDB("Unexpected error when reading events:", err, message.Chat.ID)
	}
	msg.ReplyMarkup = keyboard

	// set state to cache
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), choosingEventState)
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
	msg := tgbotapi.NewMessage(message.Chat.ID, replyText)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return msg
}
