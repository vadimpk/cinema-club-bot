package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"log"
	"net/http"
	"os"
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

		adminUpdates, _ = b.adminBot.GetUpdatesChan(u)
	} else {
		_, err := b.adminBot.SetWebhook(tgbotapi.NewWebhook(fmt.Sprintf(cfg.Heroku.URL, b.adminBot.Token)))
		if err != nil {
			return err
		}

		adminUpdates = b.adminBot.ListenForWebhook("/" + b.adminBot.Token)
	}

	// if debug - polling
	if cfg.PublicBot.Debug {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = cfg.PublicBot.Timeout

		publicUpdates, _ = b.publicBot.GetUpdatesChan(u)
	} else {
		_, err := b.publicBot.SetWebhook(tgbotapi.NewWebhook(fmt.Sprintf(cfg.Heroku.URL, b.publicBot.Token)))
		if err != nil {
			return err
		}

		publicUpdates = b.publicBot.ListenForWebhook("/" + b.publicBot.Token)
	}

	go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			log.Println(err)
		}
	}()

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
