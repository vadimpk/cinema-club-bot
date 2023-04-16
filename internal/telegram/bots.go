package telegram

import (
	"github.com/vadimpk/cinema-club-bot/config"
	"github.com/vadimpk/cinema-club-bot/pkg/event"
	"log"
	"net/http"
)

type Bots struct {
	adminBot  *Bot
	publicBot *Bot
	parseMode string
	eventBus  event.Bus
}

func NewBots(adminBot *Bot, publicBot *Bot, bus event.Bus) *Bots {
	bots := &Bots{adminBot: adminBot, publicBot: publicBot, eventBus: bus}

	err := bus.Subscribe("send messages from admin", bots.publicBot.sendMessagesFromAdmin)
	if err != nil {
		log.Fatal("failed to subscribe to event created event", "error", err)
	}

	return bots
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
