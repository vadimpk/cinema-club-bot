package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/configs"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"github.com/vadimpk/cinema-club-bot/internal/handlers"
	"github.com/vadimpk/cinema-club-bot/internal/repository"
	"github.com/vadimpk/cinema-club-bot/pkg/logging"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	handler    handlers.Handler
	cache      cache.Cache
	repository repository.Repositories
	updates    tgbotapi.UpdatesChannel
	parseMode  string
	logger     logging.Logger
}

func NewBot(bot *tgbotapi.BotAPI, handler handlers.Handler, cache cache.Cache, repository repository.Repositories, logger logging.Logger) *Bot {
	return &Bot{bot: bot, handler: handler, cache: cache, repository: repository, logger: logger}
}

func Init(cfg configs.BotConfig, handler handlers.Handler, cache cache.Cache, repository repository.Repositories, logger logging.Logger) (*Bot, error) {

	bot, err := tgbotapi.NewBotAPI(cfg.TOKEN)
	if err != nil {
		return nil, err
	}

	bot.Debug = cfg.Debug

	telegramBot := NewBot(bot, handler, cache, repository, logger)
	telegramBot.SetParseMode(cfg.ParseMode)

	return telegramBot, nil
}

func (b *Bot) SetParseMode(parseMode string) {
	b.parseMode = parseMode
}

func (b *Bot) initUpdatesChannel(cfg configs.BotConfig) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Timeout

	upd, err := b.bot.GetUpdatesChan(u)
	b.updates = upd

	return err
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
