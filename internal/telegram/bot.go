package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"github.com/vadimpk/cinema-club-bot/internal/config"
	"github.com/vadimpk/cinema-club-bot/internal/handlers"
	"github.com/vadimpk/cinema-club-bot/internal/repository"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	handler    handlers.Handler
	cache      cache.Cache
	repository repository.Repositories
	updates    tgbotapi.UpdatesChannel
	parseMode  string
}

func NewBot(bot *tgbotapi.BotAPI, handler handlers.Handler, cache cache.Cache, repository repository.Repositories) *Bot {
	return &Bot{bot: bot, handler: handler, cache: cache, repository: repository}
}

func Init(cfg config.BotConfig, handler handlers.Handler, cache cache.Cache, repository repository.Repositories) (*Bot, error) {

	bot, err := tgbotapi.NewBotAPI(cfg.TOKEN)
	if err != nil {
		return nil, err
	}

	bot.Debug = cfg.Debug

	telegramBot := NewBot(bot, handler, cache, repository)
	telegramBot.SetParseMode(cfg.ParseMode)

	return telegramBot, nil
}

func (b *Bot) SetParseMode(parseMode string) {
	b.parseMode = parseMode
}

func (b *Bot) initUpdatesChannel(cfg config.BotConfig, herokuConfig config.HerokuConfig) error {
	// if debug - polling
	if cfg.Debug {
		_, _ = b.bot.SetWebhook(tgbotapi.NewWebhook(""))

		u := tgbotapi.NewUpdate(0)
		u.Timeout = cfg.Timeout

		upd, err := b.bot.GetUpdatesChan(u)
		b.updates = upd

		return err
	} else {
		// set heroku webhook
		_, err := b.bot.SetWebhook(tgbotapi.NewWebhook(fmt.Sprintf(herokuConfig.URL, b.bot.Token)))
		if err != nil {
			return err
		}

		b.updates = b.bot.ListenForWebhook("/" + b.bot.Token)
		return nil
	}
}

func (b *Bot) handleUpdates() {
	for update := range b.updates {
		if update.Message != nil { // If we got a message
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			messages := b.handler.HandleMessage(update.Message)
			for _, msg := range messages {
				b.sendMessage(msg)
			}
		}
	}
}
