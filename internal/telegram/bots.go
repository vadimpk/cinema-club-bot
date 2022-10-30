package telegram

import (
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
	err := b.adminBot.initUpdatesChannel(cfg.AdminBot)
	if err != nil {
		return err
	}
	err = b.publicBot.initUpdatesChannel(cfg.PublicBot)
	if err != nil {
		return err
	}

	go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			log.Println(err)
		}
	}()

	go b.adminBot.handleUpdates()
	b.publicBot.handleUpdates()
	return nil
}
