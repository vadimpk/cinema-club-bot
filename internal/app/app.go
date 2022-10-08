package app

import (
	redis2 "github.com/vadimpk/cinema-club-bot/internal/cache/redis"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/admin"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/public"
	"github.com/vadimpk/cinema-club-bot/internal/telegram"
	"log"
)

func Run(configDir string) {

	cfg, err := config.Init(configDir)
	if err != nil {
		log.Fatal(err)
	}

	cache := redis2.NewCache(cfg.Redis)

	adminHandler := admin.NewHandler(cache)
	publicHandler := public.NewHandler(cache)

	adminBot, err := telegram.Init(cfg.AdminBot, adminHandler, cache)
	if err != nil {
		log.Fatal(err)
	}

	publicBot, err := telegram.Init(cfg.PublicBot, publicHandler, cache)
	if err != nil {
		log.Fatal(err)
	}

	bots := telegram.NewBots(adminBot, publicBot)

	if err := bots.Start(cfg); err != nil {
		log.Fatal(err)
	}

}
