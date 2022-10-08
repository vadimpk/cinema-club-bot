package config

import (
	"github.com/spf13/viper"
	"os"
)

type (
	Config struct {
		AdminBot  BotConfig
		PublicBot BotConfig
		Heroku    HerokuConfig
		Redis     RedisConfig
	}

	BotConfig struct {
		Debug     bool   `mapstructure:"debug"`
		Timeout   int    `mapstructure:"timeout"`
		ParseMode string `mapstructure:"parse_mode"`
		TOKEN     string
	}

	RedisConfig struct {
		URL      string `mapstructure:"url"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}

	HerokuConfig struct {
		URL string `mapstructure:"url"`
	}
)

func Init(configPath string) (*Config, error) {
	if err := parseConfigDir(configPath, os.Getenv("APP_ENV")); err != nil {
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
	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigDir(dir, env string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// set additional config file based on environment (local / production)
	if env == "" {
		viper.SetConfigName("local")
	} else {
		viper.SetConfigName(env)
	}

	return viper.MergeInConfig()
}

func parseEnv(cfg *Config) error {

	if cfg.AdminBot.Debug {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		cfg.AdminBot.TOKEN = viper.GetString("TEST_ADMIN_BOT_API_TOKEN")
	} else {
		cfg.AdminBot.TOKEN = os.Getenv("ADMIN_BOT_API_TOKEN")
	}

	if cfg.PublicBot.Debug {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		cfg.PublicBot.TOKEN = viper.GetString("TEST_PUBLIC_BOT_API_TOKEN")
	} else {
		cfg.PublicBot.TOKEN = os.Getenv("PUBLIC_BOT_API_TOKEN")
	}

	return nil
}
