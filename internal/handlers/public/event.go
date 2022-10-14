package public

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
	"strconv"
)

func (h *Handler) seeProgram(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	// set state to cache
	err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), lookingAtProgramState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// get events
	events, err := h.repos.GetActive(ctx)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "На даний час активних подій немає")}
		}
		return h.errorDB("Unexpected error when reading active events:", err, message.Chat.ID)
	}
	text := "Найближчі події Кіноклубу:\n\n"
	for _, event := range events {
		list, err := h.repos.GetList(ctx, event.ListID)
		if err != nil {
			return h.errorDB("Unexpected error when reading active events:", err, message.Chat.ID)
		}
		text += event.PreviewForProgram(list)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = h.getToMainMenuKeyboard(false)
	return []tgbotapi.MessageConfig{msg}
}

/*
seeEvent - validates given message from previous step and if valid sets it to cache.
Returns message with keyboard to register at the event
*/
func (h *Handler) seeEvent(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
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
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	// set state to cache
	err = h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), lookingAtEventState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// set identifier to cache
	err = h.cache.SetIdentifier(ctx, convertChatIDToString(message.Chat.ID), identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, event.Preview(list))
	msg.ReplyMarkup = h.getRegisterKeyboard(false, list, convertChatIDToString(message.Chat.ID))
	return []tgbotapi.MessageConfig{msg}
}

func (h *Handler) registerAtEvent(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	// get identifier from cache
	identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// get event
	event, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	// get list
	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	// check if user is not yet registered
	for _, u := range list.List {
		if u.ChatID == convertChatIDToString(message.Chat.ID) {
			return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Ви вже зареєстровані на цю подію")}
		}
	}

	// check if enough free seats
	if len(list.List) >= list.Capacity {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Місць на подію вже немає. Перегляньте інші події в афіші")}
	}

	return h.askToEnterData(ctx, message, enteringNameState, "Введіть призвіще імʼя")
}

func (h *Handler) getName(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	// set name to cache
	err := h.cache.SetName(ctx, convertChatIDToString(message.Chat.ID), message.Text)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	return h.askToEnterData(ctx, message, enteringPhoneNumberState, "Введіть номер телефону (0681234567):")
}

func (h *Handler) getPhoneNumber(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {

	phone := message.Text
	if len(phone) != 10 || phone == "" || phone[0] != '0' {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Неправильний формат номеру")}
	}

	// get name from cache
	name, err := h.cache.GetName(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// get name from cache
	identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// get event
	event, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// get list
	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	// check if enough free seats
	if len(list.List) >= list.Capacity {
		return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Місць на подію вже немає. Перегляньте інші події в афіші")}
	}

	list.List = append(list.List, domain.User{
		Name:     name,
		Phone:    phone,
		ChatID:   convertChatIDToString(message.Chat.ID),
		UserID:   strconv.Itoa(message.From.ID),
		Username: message.From.UserName,
	})

	err = h.repos.UpdateList(ctx, list)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	return h.goToMainMenu(ctx, message, "Успішно зареєстровано.")
}

func (h *Handler) unregisterAtEvent(ctx context.Context, message *tgbotapi.Message) []tgbotapi.MessageConfig {
	// get name from cache
	identifier, err := h.cache.GetIdentifier(ctx, convertChatIDToString(message.Chat.ID))
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// get event
	event, err := h.repos.GetEvent(ctx, identifier)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}

	// get list
	list, err := h.repos.GetList(ctx, event.ListID)
	if err != nil {
		return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
	}

	// check if user is not yet registered
	for i, u := range list.List {
		if u.ChatID == convertChatIDToString(message.Chat.ID) {

			// remove value at index
			list.List[i] = list.List[len(list.List)-1]
			list.List = list.List[:len(list.List)-1]

			// update list
			err = h.repos.UpdateList(ctx, list)
			if err != nil {
				return h.errorDB("Unexpected error when getting event: ", err, message.Chat.ID)
			}
			return h.goToMainMenu(ctx, message, "Реєстрацію скасовано")
		}
	}

	return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Ви ще не зареєстровані на цю подію")}
}
