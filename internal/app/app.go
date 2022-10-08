package app

import (
	"github.com/go-redis/redis/v9"
	redis2 "github.com/vadimpk/cinema-club-bot/internal/cache/redis"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/admin"
	"github.com/vadimpk/cinema-club-bot/internal/handlers/public"
	"github.com/vadimpk/cinema-club-bot/internal/repository/mongodb"
	"github.com/vadimpk/cinema-club-bot/internal/telegram"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type PromoCode struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code string             `json:"code" bson:"code"`
}

func Run(configDir string) {

	cfg, err := config.Init(configDir)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	cache := redis2.NewCache(rdb, cfg.Redis.TTL)

	mongoClient, err := mongodb.NewClient(cfg.Mongo)
	if err != nil {
		log.Fatal(err)
	}

	_ = mongoClient.Database(cfg.Mongo.Name)

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
