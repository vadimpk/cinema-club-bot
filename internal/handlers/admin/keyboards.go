package admin

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
)

const (
	updateEventOption = "Редагувати подію"
	lookUpListsOption = "Переглянути списки"
	createEventOption = "Створити подію"

	updateEventNameOption         = "Редагувати назву"
	updateEventDescriptionOption  = "Редагувати опис"
	updateEventDateOption         = "Редагувати дату"
	updateEventListCapacityOption = "Редагувати кількість місць"
	activateEventOption           = "Активувати подію"
	deactivateEventOption         = "Деактивувати подію"
	deleteEventOption             = "Видалити подію"

	cancelUpdateOption = "Скасувати"

	deleteReservationOption = "Видалити учасника"
	sendMessageToAllOption  = "Написати всім"

	toMainMenuOption = "На головну"
)

func (h *Handler) getOptionsKeyboard(oneTime bool) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(updateEventOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(lookUpListsOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(createEventOption)))

	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}

func (h *Handler) getUpdateEventOptionsKeyboard(event domain.Event, oneTime bool) tgbotapi.ReplyKeyboardMarkup {
	var activationOption string
	if event.Active {
		activationOption = deactivateEventOption
	} else {
		activationOption = activateEventOption
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(updateEventNameOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(updateEventDescriptionOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(updateEventDateOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(updateEventListCapacityOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(activationOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(deleteEventOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(toMainMenuOption)))
	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}

func (h *Handler) getEventsKeyboard(ctx context.Context, oneTime bool) (tgbotapi.ReplyKeyboardMarkup, error) {
	events, err := h.repos.GetAll(ctx)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	buttons := make([][]tgbotapi.KeyboardButton, len(events)+1)
	for id, event := range events {
		buttons[id] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(event.Identifier))
	}
	buttons[len(events)] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(toMainMenuOption))
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	keyboard.OneTimeKeyboard = oneTime
	return keyboard, nil
}

func (h *Handler) getCancelUpdateKeyboard(oneTime bool) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(cancelUpdateOption)))
	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}

func (h *Handler) getListOptions(oneTime bool) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(deleteReservationOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(sendMessageToAllOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(toMainMenuOption)))

	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}

func (h *Handler) getToMainMenuOption(oneTime bool) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(toMainMenuOption)))

	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}
