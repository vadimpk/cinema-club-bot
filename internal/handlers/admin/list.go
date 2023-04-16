package admin

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

func (h *Handler) updateList(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(list *domain.List, text string) error, replyText string) []tgbotapi.MessageConfig {

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

	// get list from db
	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when reading event:", err, message.Chat.ID)
	}

	err = updateFunc(&list, message.Text)
	if err != nil {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, err.Error())}
	}

	// update entry in db
	err = h.repos.UpdateList(ctx, list)
	if err != nil {
		return h.errorDB("Unexpected error when updating event:", err, message.Chat.ID)
	}

	return h.goBackToUpdateOptions(ctx, message, replyText)
}

func (h *Handler) updateListOnCreation(ctx context.Context, message *tgbotapi.Message,
	updateFunc func(list *domain.List, text string) error,
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
		return h.goToMainMenu(ctx, message, "Подію не збережено скасовано.")
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

	// get list from db
	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when reading event:", err, message.Chat.ID)
	}

	err = updateFunc(&list, message.Text)
	if err != nil {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, err.Error())}
	}

	// update entry in db
	err = h.repos.UpdateList(ctx, list)
	if err != nil {
		return h.errorDB("Unexpected error when updating event:", err, message.Chat.ID)
	}
	if last {
		return h.goToMainMenu(ctx, message, replyText)
	}
	return h.askToEnterData(ctx, message, state, replyText)
}

func (h *Handler) updateListCapacity(list *domain.List, text string) error {
	i, err := strconv.Atoi(text)
	if err != nil {
		return errors.New("Неправильний формат даних.")
	}
	if i < len(list.List) {
		return errors.New("Не можна поставити менше число, ніж кількість зареєстрованих учасників на даний момент. Спробуйте видалити деяких учасників або поставити більше число.")
	}
	list.Capacity = i
	return nil
}

func (h *Handler) retrieveList(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
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
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Події з таким ідентифікатором не існує. Виберіть ідентифікатор ще раз:")}
	}

	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when getting list: ", err, message.Chat.ID)
	}

	// set state to cache
	err = h.cache.SetAdminState(ctx, convertChatIDToString(message.Chat.ID), lookingAtListState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// set identifier to cache
	err = h.cache.SetAdminIdentifier(ctx, convertChatIDToString(message.Chat.ID), identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, list.Preview())
	msg.ReplyMarkup = h.getListOptions(true)
	return []tgbotapi.MessageConfig{msg}
}

func (h *Handler) deleteReservation(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	if message.Text == cancelUpdateOption {
		return h.goToMainMenu(ctx, message, "Скасовано")
	}
	// get state from cache
	identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	event, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when getting list: ", err, message.Chat.ID)
	}

	for i, u := range list.List {
		if u.Phone == message.Text {
			// remove value at index
			list.List[i] = list.List[len(list.List)-1]
			list.List = list.List[:len(list.List)-1]

			// update list
			err = h.repos.UpdateList(ctx, list)
			if err != nil {
				return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
			}

			messages := []domain.Message{{
				ChatID: u.ChatID,
				Text:   "Вас було видалено з події " + list.EventIdentifier,
			}}

			h.eventBus.Publish("send messages from admin", messages)
			return h.goToMainMenu(ctx, message, fmt.Sprintf("Реєстрацію користувача %s видалено", u.Name))
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Користувача з таким телефоном не знайдено, спробуйте ще раз")
	msg.ReplyMarkup = h.getCancelUpdateKeyboard(true)
	return []tgbotapi.MessageConfig{msg}
}

func (h *Handler) sendMessageToAll(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {

	if message.Text == cancelUpdateOption {
		return h.goToMainMenu(ctx, message, "Скасовано")
	}

	// get identifier from cache
	identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	event, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when getting list: ", err, message.Chat.ID)
	}
	messagesToSend := make([]domain.Message, len(list.List))
	for i, u := range list.List {
		messagesToSend[i] = domain.Message{
			ChatID: u.ChatID,
			Text:   message.Text,
		}
	}

	h.eventBus.Publish("send messages from admin", messagesToSend)

	return h.goToMainMenu(ctx, message, fmt.Sprintf("%d повідомлень буде надіслано в іншому боті", len(messagesToSend)))
}
