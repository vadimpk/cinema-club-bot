package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vadimpk/cinema-club-bot/internal/config"

	"log"
)

type Bot struct {
	adminBot  *tgbotapi.BotAPI
	publicBot *tgbotapi.BotAPI
	parseMode string
}

func NewBot(adminBot *tgbotapi.BotAPI, publicBot *tgbotapi.BotAPI) *Bot {
	return &Bot{adminBot: adminBot, publicBot: publicBot}
}

func (b *Bot) SetParseMode(parseMode string) {
	b.parseMode = parseMode
}

func (b *Bot) Start(cfg *config.Config) error {

	// init
	adminUpdates := tgbotapi.UpdatesChannel(make(chan tgbotapi.Update))
	publicUpdates := tgbotapi.UpdatesChannel(make(chan tgbotapi.Update))

	// if debug - polling
	if cfg.AdminBot.Debug {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = cfg.AdminBot.Timeout

		adminUpdates = b.adminBot.GetUpdatesChan(u)
	} else {
		// set webhooks (heroku)
	}

	// if debug - polling
	if cfg.PublicBot.Debug {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = cfg.PublicBot.Timeout

		publicUpdates = b.publicBot.GetUpdatesChan(u)
	} else {
		// set webhooks (heroku)
	}

	go b.handleUpdates(adminUpdates, b.adminBot)
	b.handleUpdates(publicUpdates, b.publicBot)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text))

			//if update.Message.IsCommand() {
			//	if err := b.handleCommands(update.Message); err != nil {
			//		b.handleError(update.Message.Chat.ID, err)
			//	}
			//	continue
			//}
			//
			//if err := b.handleMessage(update.Message); err != nil {
			//	b.handleError(update.Message.Chat.ID, err)
			//}
		}
	}
}
