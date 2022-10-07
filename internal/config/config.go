package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

type (
	Config struct {
		AdminBot  BotConfig
		PublicBot BotConfig
		Heroku    HerokuConfig
	}

	BotConfig struct {
		Debug     bool   `mapstructure:"debug"`
		Timeout   int    `mapstructure:"timeout"`
		ParseMode string `mapstructure:"parse_mode"`
		TOKEN     string
	}

	HerokuConfig struct {
		URL string `mapstructure:"url"`
	}
)

func Init(configPath string) (*Config, error) {
	if err := parseConfigPath(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("admin-bot", &cfg.AdminBot); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("public-bot", &cfg.PublicBot); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("heroku", &cfg.Heroku); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigPath(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}

func parseEnv(cfg *Config) error {

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if cfg.AdminBot.Debug {
		cfg.AdminBot.TOKEN = viper.GetString("TEST_ADMIN_BOT_API_TOKEN")
	} else {
		cfg.AdminBot.TOKEN = os.Getenv("ADMIN_BOT_API_TOKEN")
	}

	if cfg.PublicBot.Debug {
		cfg.PublicBot.TOKEN = viper.GetString("TEST_PUBLIC_BOT_API_TOKEN")
	} else {
		cfg.PublicBot.TOKEN = os.Getenv("PUBLIC_BOT_API_TOKEN")
	}

	return nil
}
