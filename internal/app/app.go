package app

import (
	"github.com/vadimpk/cinema-club-bot/configs"
	"github.com/vadimpk/cinema-club-bot/internal/cache/in_memory"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/admin"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/public"
	"github.com/vadimpk/cinema-club-bot/internal/repository/mongodb"
	"github.com/vadimpk/cinema-club-bot/internal/telegram"
	"github.com/vadimpk/cinema-club-bot/pkg/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromoCode struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code string             `json:"code" bson:"code"`
}

func Run(configDir, configsFile string) {

	logger := logging.NewZap("INFO")

	cfg, err := configs.Init(configDir, configsFile)
	if err != nil {
		logger.Fatal("failed to init configs", "error", err)
	}

	cache := in_memory.NewCache(cfg.Cache.TTL, cfg.Cache.AdminTTL)

	mongoClient, err := mongodb.NewClient(cfg.Mongo)
	if err != nil {
		logger.Fatal("failed to init mongo client", "error", err)
	}

	mdb := mongoClient.Database(cfg.Mongo.Name)

	repos := mongodb.NewRepositories(mdb, logger.Named("mongo"))

	adminHandler := admin.NewHandler(cache, repos, logger.Named("admin_handler"))
	publicHandler := public.NewHandler(cache, repos, logger.Named("public_handler"))

	adminBot, err := telegram.Init(cfg.AdminBot, adminHandler, cache, repos, logger.Named("admin_bot"))
	if err != nil {
		logger.Fatal("failed to init admin bot", "error", err)
	}

	publicBot, err := telegram.Init(cfg.PublicBot, publicHandler, cache, repos, logger.Named("public_bot"))
	if err != nil {
		logger.Fatal("failed to init public bot", "error", err)
	}

	bots := telegram.NewBots(adminBot, publicBot)

	if err := bots.Start(cfg); err != nil {
		logger.Fatal("failed to start bots", "error", err)
	}
}
