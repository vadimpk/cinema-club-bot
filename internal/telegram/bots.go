package telegram

import (
	"github.com/vadimpk/cinema-club-bot/configs"
	"log"
	"net/http"
)

type Bots struct {
	adminBot  *Bot
	publicBot *Bot
	parseMode string
}

func NewBots(adminBot *Bot, publicBot *Bot) *Bots {
	return &Bots{adminBot: adminBot, publicBot: publicBot}
}

func (b *Bots) Start(cfg *configs.Config) error {

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
		log.Println("starting http server on port " + cfg.HTTP.Port)
		err := http.ListenAndServe(cfg.HTTP.Port, nil)
		if err != nil {
			log.Println(err)
		}
	}()

	go b.adminBot.handleUpdates()
	b.publicBot.handleUpdates()
	return nil
}
