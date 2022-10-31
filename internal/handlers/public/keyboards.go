package public

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
)

func (h *Handler) getOptionsKeyboard(oneTime bool) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(seeProgramOption)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(registerAtEventOption)))

	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}

func (h *Handler) getToMainMenuKeyboard(oneTime bool) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(toMainMenuOption)))
	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}

func (h *Handler) getEventsKeyboard(ctx context.Context, oneTime bool) (tgbotapi.ReplyKeyboardMarkup, error) {
	events, err := h.repos.GetActive(ctx)
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

func (h *Handler) getRegisterKeyboard(oneTime bool, list domain.List, chatID string) tgbotapi.ReplyKeyboardMarkup {
	option := registerOption
	if len(list.List) >= list.Capacity {
		option = noSeatsOption
	}
	for _, u := range list.List {
		if u.ChatID == chatID {
			option = unregisterOption
		}
	}
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(option)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(toMainMenuOption)))
	keyboard.OneTimeKeyboard = oneTime
	return keyboard
}
