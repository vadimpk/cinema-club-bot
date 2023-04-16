package app

import (
	"github.com/vadimpk/cinema-club-bot/config"
	"github.com/vadimpk/cinema-club-bot/internal/cache/in_memory"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/admin"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/public"
	"github.com/vadimpk/cinema-club-bot/internal/repository/mongodb"
	"github.com/vadimpk/cinema-club-bot/internal/telegram"
	"github.com/vadimpk/cinema-club-bot/pkg/event"
	"github.com/vadimpk/cinema-club-bot/pkg/logging"
)

func Run() {

	logger := logging.NewZap("info")
	cfg := config.Get()
	logger.Info("got config", "cfg", cfg)

	cache := in_memory.NewCache(cfg.Cache.TTL, cfg.Cache.AdminTTL)

	mongoClient, err := mongodb.NewClient(cfg.Mongo)
	if err != nil {
		logger.Fatal("failed to init mongo client", "error", err)
	}

	eventBus := event.NewMsgBus(128)

	mdb := mongoClient.Database(cfg.Mongo.DBName)

	repos := mongodb.NewRepositories(mdb, logger.Named("mongo"))

	adminHandler := admin.NewHandler(cache, repos, logger.Named("admin_handler"), eventBus)
	publicHandler := public.NewHandler(cache, repos, logger.Named("public_handler"))

	adminBot, err := telegram.Init(cfg.AdminBot, adminHandler, cache, repos, logger.Named("admin_bot"))
	if err != nil {
		logger.Fatal("failed to init admin bot", "error", err)
	}

	publicBot, err := telegram.Init(cfg.PublicBot, publicHandler, cache, repos, logger.Named("public_bot"))
	if err != nil {
		logger.Fatal("failed to init public bot", "error", err)
	}

	bots := telegram.NewBots(adminBot, publicBot, eventBus)

	if err := bots.Start(cfg); err != nil {
		logger.Fatal("failed to start bots", "error", err)
	}
}
