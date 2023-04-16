package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

var (
	errMessageIsTooLong = "Bad Request: message is too long"
)

func (b *Bot) sendMessage(msg tgbotapi.MessageConfig) {
	msg.ParseMode = b.parseMode
	if msg.ParseMode == "markdown" {
		msg.Text = replaceReservedCharacters(msg.Text)
	}
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
			b.logger = b.logger.With("text", msg.Text)
			_, err := b.bot.Send(msg)
			if err != nil {
				b.logger.Error("failed to send message", err)
			}
		}
	}
}

func replaceReservedCharacters(text string) string {
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "*", "\\*")
	return text
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
