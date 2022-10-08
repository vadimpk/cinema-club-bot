package app

import (
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/admin"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/public"
	"github.com/vadimpk/cinema-club-bot/internal/telegram"
	"log"
)

func Run(configPath string) {

	cfg, err := config.Init(configPath)

	if err != nil {
		log.Fatal(err)
	}

	adminHandler := admin.NewHandler()
	publicHandler := public.NewHandler()

	adminBot, err := telegram.Init(cfg.AdminBot, adminHandler)
	if err != nil {
		log.Fatal(err)
	}

	publicBot, err := telegram.Init(cfg.PublicBot, publicHandler)
	if err != nil {
		log.Fatal(err)
	}

	bots := telegram.NewBots(adminBot, publicBot)

	if err := bots.Start(cfg); err != nil {
		log.Fatal(err)
	}

}
