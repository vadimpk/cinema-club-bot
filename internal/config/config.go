package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

type (
	Config struct {
		Main      MainConfig
		AdminBot  BotConfig
		PublicBot BotConfig
		Web       WebConfig
		Redis     RedisConfig
		Mongo     MongoConfig
	}

	MainConfig struct {
		Environment string
		Stage       string
	}

	BotConfig struct {
		Debug     bool   `mapstructure:"debug"`
		Timeout   int    `mapstructure:"timeout"`
		ParseMode string `mapstructure:"parse_mode"`
		TOKEN     string
	}

	RedisConfig struct {
		URL      string        `mapstructure:"url"`
		Port     string        `mapstructure:"port"`
		Password string        `mapstructure:"password"`
		DB       int           `mapstructure:"db"`
		TTL      time.Duration `mapstructure:"ttl"`
	}

	MongoConfig struct {
		URI  string
		Name string
	}

	WebConfig struct {
		URL string
	}
)

func Init(configPath, configsFile string) (*Config, error) {

	if err := parseConfigDir(configPath, configsFile); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("main", &cfg.Main); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("admin-bot", &cfg.AdminBot); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("public-bot", &cfg.PublicBot); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigDir(dir, file string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)

	return viper.ReadInConfig()
}

func parseEnv(cfg *Config) error {

	if cfg.Main.Environment == "local" {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		cfg.Mongo.URI = viper.GetString("MONGO_URI")
		cfg.Mongo.Name = viper.GetString("MONGO_DB_NAME")
		cfg.AdminBot.TOKEN = viper.GetString("TEST_ADMIN_BOT_API_TOKEN")
		cfg.PublicBot.TOKEN = viper.GetString("TEST_PUBLIC_BOT_API_TOKEN")

	} else {
		cfg.Mongo.URI = os.Getenv("MONGO_URI")
		cfg.Mongo.Name = os.Getenv("MONGO_DB_NAME")

		cfg.AdminBot.TOKEN = os.Getenv("ADMIN_BOT_API_TOKEN")
		cfg.PublicBot.TOKEN = os.Getenv("PUBLIC_BOT_API_TOKEN")

		cfg.Web.URL = os.Getenv("AUTH_SERVER_URL")
	}

	return nil
}
