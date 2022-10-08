package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"github.com/vadimpk/cinema-club-bot/internal/handlers"
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	handler   handlers.Handler
	parseMode string
}

func NewBot(bot *tgbotapi.BotAPI, handler handlers.Handler) *Bot {
	return &Bot{bot: bot, handler: handler}
}

func Init(cfg config.BotConfig, handler handlers.Handler) (*Bot, error) {

	bot, err := tgbotapi.NewBotAPI(cfg.TOKEN)
	if err != nil {
		return nil, err
	}

	bot.Debug = cfg.Debug

	telegramBot := NewBot(bot, handler)
	telegramBot.SetParseMode(cfg.ParseMode)

	return telegramBot, nil
}

func (b *Bot) SetParseMode(parseMode string) {
	b.parseMode = parseMode
}

func (b *Bot) initUpdatesChannel(cfg config.BotConfig, herokuConfig config.HerokuConfig) (tgbotapi.UpdatesChannel, error) {
	// if debug - polling
	if cfg.Debug {
		_, _ = b.bot.SetWebhook(tgbotapi.NewWebhook(""))

		u := tgbotapi.NewUpdate(0)
		u.Timeout = cfg.Timeout

		return b.bot.GetUpdatesChan(u)
	} else {
		// set heroku webhook
		_, err := b.bot.SetWebhook(tgbotapi.NewWebhook(fmt.Sprintf(herokuConfig.URL, b.bot.Token)))
		if err != nil {
			return nil, err
		}

		return b.bot.ListenForWebhook("/" + b.bot.Token), nil
	}
}
