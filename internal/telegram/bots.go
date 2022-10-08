package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"log"
	"net/http"
	"os"
)

type Bots struct {
	adminBot  *Bot
	publicBot *Bot
	parseMode string
}

func NewBots(adminBot *Bot, publicBot *Bot) *Bots {
	return &Bots{adminBot: adminBot, publicBot: publicBot}
}

func (b *Bots) Start(cfg *config.Config) error {

	// init
	adminUpdates, err := b.adminBot.initUpdatesChannel(cfg.AdminBot, cfg.Heroku)
	if err != nil {
		return err
	}
	publicUpdates, err := b.publicBot.initUpdatesChannel(cfg.PublicBot, cfg.Heroku)
	if err != nil {
		return err
	}

	go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			log.Println(err)
		}
	}()

	go b.handleUpdates(adminUpdates, b.adminBot.bot)
	b.handleUpdates(publicUpdates, b.publicBot.bot)
	return nil
}

func (b *Bots) handleUpdates(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch bot.Token {
			case b.adminBot.bot.Token:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "from admin bot"+update.Message.Text))
			case b.publicBot.bot.Token:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "from public bot"+update.Message.Text))
			}
		}
	}
}
