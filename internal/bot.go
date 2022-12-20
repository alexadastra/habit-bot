package internal

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
    botAPI *tgbotapi.BotAPI
}

func NewBot(botAPI *tgbotapi.BotAPI) *Bot {
    return &Bot{
        botAPI: botAPI,
    }
}

func (b *Bot) SendMessage(chatID int64, text string) error {
    msg := tgbotapi.NewMessage(chatID, text)
    _, err := b.botAPI.Send(msg)
    return err
}
