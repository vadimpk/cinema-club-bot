package admin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (h *Handler) errorDB(msg string, err error, chatID int64) (tgbotapi.MessageConfig, error) {
	log.Println(msg, err)
	return tgbotapi.NewMessage(chatID, "Сталася помилка. Натисніть команду /start"), err
}
