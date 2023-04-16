package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/cinema-club-bot/config"
	"github.com/vadimpk/cinema-club-bot/internal/cache"
	"github.com/vadimpk/cinema-club-bot/internal/domain"
	"github.com/vadimpk/cinema-club-bot/internal/handlers"
	"github.com/vadimpk/cinema-club-bot/internal/repository"
	"github.com/vadimpk/cinema-club-bot/pkg/logging"
	"strconv"
	"strings"
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

func Init(cfg config.BotConfig, handler handlers.Handler, cache cache.Cache, repository repository.Repositories, logger logging.Logger) (*Bot, error) {

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
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

func (b *Bot) initUpdatesChannel(cfg config.BotConfig) error {
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

func (b *Bot) sendMessagesFromAdmin(msgs []domain.Message) {
	messagesToSend := make([]tgbotapi.MessageConfig, 0)
	for _, m := range msgs {
		chatID, err := strconv.Atoi(m.ChatID)
		if err == nil {
			messagesToSend = append(messagesToSend, tgbotapi.NewMessage(int64(chatID), m.Text))
		}
	}

	for _, msg := range messagesToSend {
		msg.Text = replaceReservedCharacters(msg.Text)
		b.sendMessage(msg)
	}
}

func replaceReservedCharacters(text string) string {
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "*", "\\*")
	return text
}
