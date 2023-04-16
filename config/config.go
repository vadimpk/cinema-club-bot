package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
	"time"
)

type (
	Config struct {
		AdminBot  BotConfig
		PublicBot BotConfig
		Mongo     MongoConfig
		Cache     CacheConfig
		HTTP      HTTPConfig
	}

	BotConfig struct {
		Debug     bool   `env:"DEFAULT_BOT_DEBUG" env-default:"true"`
		Timeout   int    `env:"DEFAULT_BOT_TIMEOUT" env-default:"60"`
		ParseMode string `env:"DEFAULT_BOT_PARSE_MODE" env-default:"markdown"`
		Token     string
	}

	MongoConfig struct {
		URI    string `env:"MONGO_URI"`
		DBName string `env:"MONGO_DB_NAME"`
	}

	HTTPConfig struct {
		Port string `env:"HTTP_PORT" env-default:":8080"`
	}

	CacheConfig struct {
		TTL      time.Duration `env:"CACHE_PORT" env-default:"4h"`
		AdminTTL time.Duration `env:"ADMIN_CACHE_TTL" env-default:"720h"`
	}
)

var (
	config Config
	once   sync.Once
)

// Get returns a config.
// Get loads the config only once.
func Get() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}

		err = cleanenv.ReadEnv(&config)
		if err != nil {
			log.Fatal("failed to read env", err)
		}

		config.AdminBot.Token = os.Getenv("ADMIN_BOT_TOKEN")
		config.PublicBot.Token = os.Getenv("PUBLIC_BOT_TOKEN")
	})

	return &config
}
