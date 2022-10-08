package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Handler interface {
	HandleMessage(message *tgbotapi.Message) (tgbotapi.MessageConfig, error)
}
