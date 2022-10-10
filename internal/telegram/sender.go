package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var (
	errMessageIsTooLong = "Bad Request: message is too long"
)

func (b *Bot) sendMessage(msg tgbotapi.MessageConfig) {
	msg.ParseMode = b.parseMode
	_, err := b.bot.Send(msg)

	if err != nil {
		switch err.Error() {
		case errMessageIsTooLong:
			splitted := split(msg.Text, len(msg.Text)/2)
			for _, s := range splitted {
				msg := tgbotapi.NewMessage(msg.ChatID, s)
				b.sendMessage(msg)
			}
		default:
			_, err := b.bot.Send(msg)
			log.Println(err)
		}
	}
}

func split(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
