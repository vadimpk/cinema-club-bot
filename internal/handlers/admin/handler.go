package admin

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"github.com/vadimpk/cinema-club-bot/internal/repository"
	"strconv"
)

type Handler struct {
	cache cache.Cache
	repos repository.Repositories
}

const (
	startState  = "starting state"
	createState = "creating new event state"

	updateNameState         = "updating name of the event"
	updateDescriptionState  = "updating description of the event"
	updateListCapacityState = "updating list capacity of the event"
	updateDateState         = "updating state of the event"

	updateNameStateOnCreation         = "updating name of the event on creation"
	updateDescriptionStateOnCreation  = "updating description of the event on creation"
	updateListCapacityStateOnCreation = "updating list capacity of the event on creation"
	updateDateStateOnCreation         = "updating state of the event on creation"

	chooseEventToUpdateState      = "choosing event to then update"
	chooseEventForListsState      = "choosing event to see lists"
	chooseEventUpdateOptionsState = "choosing event update options"

	lookingAtListState     = "looking at list"
	deleteReservationState = "deleting reservation"
)

func NewHandler(cache cache.Cache, repos repository.Repositories) *Handler {
	return &Handler{cache: cache, repos: repos}
}

func (h *Handler) HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var ctx = context.Background()
	chatID := convertChatIDToString(message.Chat.ID)

	state, err := h.cache.GetState(ctx, chatID)
	if err != nil {
		// if state not found try to init admin
		if err == redis.Nil {
			err := h.initAdmin(ctx, chatID)
			// if unable to init admin, return message saying forbidden
			if err != nil {
				return tgbotapi.NewMessage(message.Chat.ID, err.Error())
			}
			state = startState
		} else {
			return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
		}
	}

	if message.IsCommand() && message.Command() == "start" {
		return h.handleStart(ctx, message, "Виберіть дію:")
	}

	if state == startState || state == chooseEventUpdateOptionsState || state == lookingAtListState {
		switch message.Text {
		case toMainMenuOption:
			return h.goToMainMenu(ctx, message, "Виберіть дію:")
		case createEventOption:
			return h.handleCreationStartProcess(ctx, message)
		case lookUpListsOption:
			return h.chooseEvent(ctx, message, chooseEventForListsState)
		case updateEventOption:
			return h.chooseEvent(ctx, message, chooseEventToUpdateState)
		case updateEventNameOption:
			return h.askToEnterData(ctx, message, updateNameState, "Введіть назву події:")
		case updateEventDescriptionOption:
			return h.askToEnterData(ctx, message, updateDescriptionState, "Введіть опис події:")
		case updateEventDateOption:
			return h.askToEnterData(ctx, message, updateDateState, "Введіть дату події (2022-02-22T15:30:00):")
		case updateEventListCapacityOption:
			return h.askToEnterData(ctx, message, updateListCapacityState, "Введіть кількість вільних місць на подію:")
		case activateEventOption:
			return h.updateEvent(ctx, message, h.updateEventActiveStatus, "Подію успішно активовано")
		case deactivateEventOption:
			return h.updateEvent(ctx, message, h.updateEventActiveStatus, "Подію успішно деактивовано")
		case deleteEventOption:
			return h.deleteEvent(ctx, message)
		case deleteReservationOption:
			return h.askToEnterData(ctx, message, deleteReservationState, "Введіть номер телефону учасника, якого видалити:")
		}
	}

	switch state {
	case startState:
		return h.handleStart(ctx, message, "Виберіть дію:")
	case chooseEventToUpdateState:
		return h.chooseUpdateOptions(ctx, message)
	case chooseEventForListsState:
		return h.retrieveList(ctx, message)
	case createState:
		return h.createNewEvent(ctx, message)
	case updateNameState:
		return h.updateEvent(ctx, message, h.updateEventName, "Назву успішно змінено")
	case updateDescriptionState:
		return h.updateEvent(ctx, message, h.updateEventDescription, "Опис успішно змінено")
	case updateListCapacityState:
		return h.updateList(ctx, message, h.updateListCapacity, "Кількість місць успішно змінено")
	case updateDateState:
		return h.updateEvent(ctx, message, h.updateEventDate, "Дату успішно змінено")
	case updateNameStateOnCreation:
		return h.updateEventOnCreation(ctx, message, h.updateEventName, updateDescriptionStateOnCreation, "Введіть опис події:", false)
	case updateDescriptionStateOnCreation:
		return h.updateEventOnCreation(ctx, message, h.updateEventDescription, updateListCapacityStateOnCreation, "Введіть кількість місць на подію:", false)
	case updateListCapacityStateOnCreation:
		return h.updateListOnCreation(ctx, message, h.updateListCapacity, updateDateStateOnCreation, "Введіть дату події (2022-02-22T15:30:00):", false)
	case updateDateStateOnCreation:
		return h.updateEventOnCreation(ctx, message, h.updateEventDate, updateDateStateOnCreation, "Подію успішно створено", true)
	case deleteReservationState:
		return h.deleteReservation(ctx, message)
	}

	return tgbotapi.MessageConfig{}
}

func (h *Handler) handleStart(ctx context.Context, message *tgbotapi.Message, msgText string) tgbotapi.MessageConfig {
	// set state to cache
	err := h.cache.SetState(ctx, convertChatIDToString(message.Chat.ID), startState)
	if err != nil {
		return h.errorDB("Unexpected error when writing cache:", err, message.Chat.ID)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = h.getOptionsKeyboard(false)
	return msg
}

func (h *Handler) initAdmin(ctx context.Context, chatID string) error {
	ok := h.repos.IsAdmin(ctx, chatID)
	if ok {
		return nil
	}
	return errors.New("forbidden")
}

func convertChatIDToString(chatID int64) string {
	return strconv.FormatInt(chatID, 10)
}
