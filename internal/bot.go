package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	// Set up the telegram bot API
	tgBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init bot")
	}
	return &Bot{
		botAPI: tgBot,
	}, nil
}

func (b *Bot) GetUpdatesChan() (tgbotapi.UpdatesChannel, error) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates, err := b.botAPI.GetUpdatesChan(updateConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updates")
	}
	return updates, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.botAPI.Send(msg)
	return err
}
