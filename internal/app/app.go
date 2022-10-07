package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"github.com/vadimpk/cinema-club-bot/internal/telegram"
	"log"
)

func Run(configPath string) {

	cfg, err := config.Init(configPath)

	if err != nil {
		log.Fatal(err)
	}

	adminBot, err := tgbotapi.NewBotAPI(cfg.AdminBot.TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	publicBot, err := tgbotapi.NewBotAPI(cfg.PublicBot.TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	adminBot.Debug = cfg.AdminBot.Debug
	publicBot.Debug = cfg.PublicBot.Debug

	telegramBot := telegram.NewBot(adminBot, publicBot)
	telegramBot.SetParseMode(cfg.AdminBot.ParseMode)

	if err := telegramBot.Start(cfg); err != nil {
		log.Fatal(err)
	}

}
