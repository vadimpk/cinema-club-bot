package public

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"github.com/vadimpk/cinema-club-bot/internal/repository"
	"github.com/vadimpk/cinema-club-bot/pkg/logging"
	"strconv"
)

type Handler struct {
	cache  cache.Cache
	repos  repository.Repositories
	logger logging.Logger
}

func NewHandler(cache cache.Cache, repos repository.Repositories, logger logging.Logger) *Handler {
	return &Handler{cache: cache, repos: repos, logger: logger}
}

const (
	startState               = "start state"
	choosingEventState       = "choosing event"
	lookingAtProgramState    = "looking at event"
	lookingAtEventState      = "looking at event"
	enteringNameState        = "entering name for registration"
	enteringPhoneNumberState = "entering phone number for registration"
)

const (
	seeProgramOption      = "Переглянути афішу"
	registerAtEventOption = "Зареєструватися на подію"
	toMainMenuOption      = "Головне меню"

	registerOption   = "Зареєструватися"
	unregisterOption = "Скасувати реєстрацію"
	noSeatsOption    = "Місць не залишилося"
)

func (h *Handler) HandleMessage(message *tgbotapi.Message) []tgbotapi.MessageConfig {
	var ctx = context.Background()
	chatID := convertChatIDToString(message.Chat.ID)

	state, err := h.cache.GetState(ctx, chatID)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}
	if state == "" {
		state = startState
	}

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			return h.handleStart(ctx, message, "Привіт!\nЦей бот створено задля перегляду афіші та реєстрації на кінопокази кіноклубу Могилянки :)\nПроблеми з реєстрацією - @vadimpk")
		}
	}

	if state == startState {
		switch message.Text {
		case seeProgramOption:
			return h.seeProgram(ctx, message)
		case registerAtEventOption:
			return h.chooseEvent(ctx, message)
		}
	}

	if state == lookingAtProgramState && message.Text == toMainMenuOption {
		return h.goToMainMenu(ctx, message, "Оберіть дію:")
	}

	if state == lookingAtEventState {
		switch message.Text {
		case registerOption:
			return h.registerAtEvent(ctx, message)
		case unregisterOption:
			return h.unregisterAtEvent(ctx, message)
		case noSeatsOption:
			return []tgbotapi.MessageConfig{tgbotapi.NewMessage(message.Chat.ID, "Місць на подію вже немає. Перегляньте інші події в афіші")}
		case toMainMenuOption:
			return h.goToMainMenu(ctx, message, "Оберіть дію:")
		}
	}

	switch state {
	case startState:
		return h.handleStart(ctx, message, "Оберіть дію:")
	case choosingEventState:
		return h.seeEvent(ctx, message)
	case enteringNameState:
		return h.getName(ctx, message)
	case enteringPhoneNumberState:
		return h.getPhoneNumber(ctx, message)
	}

	return []tgbotapi.MessageConfig{}
}

func (h *Handler) handleStart(ctx context.Context, message *tgbotapi.Message, msgText string) []tgbotapi.MessageConfig {
	// set state to cache
	err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), startState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = h.getOptionsKeyboard(false)
	return []tgbotapi.MessageConfig{msg}
}

func convertChatIDToString(chatID int64) string {
	return strconv.FormatInt(chatID, 10)
}
