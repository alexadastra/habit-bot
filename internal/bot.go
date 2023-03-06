package internal

import (
	"context"
	"log"
	"sync"

	"github.com/alexadastra/habit_bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
	inCh   tgbotapi.UpdatesChannel

	cancel context.CancelFunc
	wg     *sync.WaitGroup

	commandCh chan models.UserCommand
	messageCh chan models.UserMessage
}

// NewBot sets up the telegram bot API
func NewBot(token string) (*Bot, error) {
	tgBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init bot")
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := tgBot.GetUpdatesChan(updateConfig)

	return &Bot{
		botAPI:    tgBot,
		inCh:      updates,
		wg:        &sync.WaitGroup{},
		commandCh: make(chan models.UserCommand),
		messageCh: make(chan models.UserMessage),
	}, nil
}

func (b *Bot) GetCommandsChan() chan models.UserCommand {
	return b.commandCh
}

func (b *Bot) GetMessagesChan() chan models.UserMessage {
	return b.messageCh
}

func (b *Bot) Start(ctx context.Context) error {
	ctx, b.cancel = context.WithCancel(ctx)
	b.wg.Add(1)
	defer b.wg.Done()
	for {
		select {
		case <-ctx.Done():
			// TODO: log warning here
			return nil
		case update, ok := <-b.inCh:
			if !ok {
				// TODO: log warning here
				return nil
			}

			msg := update.Message
			if msg == nil {
				continue
			}

			// Handle commands
			if msg.IsCommand() {
				b.commandCh <- newDomainCommand(msg)
			}

			// Handle messages that are not commands
			b.messageCh <- newDomainMessage(update.Message)
		}
	}
}

func (b *Bot) Stop() {
	b.cancel()
	close(b.messageCh)
	close(b.commandCh)
	b.wg.Wait()
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	resp, err := b.botAPI.Request(tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "checkin",
			Description: "Check-in",
		},
		tgbotapi.BotCommand{
			Command:     "gratitude",
			Description: "Add gratitude",
		},
	))

	log.Println(*resp)

	if err != nil {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, err = b.botAPI.Send(msg)
	return err
}

func newDomainMessage(msg *tgbotapi.Message) models.UserMessage {
	return models.UserMessage{
		UserID:  msg.From.ID,
		Message: msg.Text,
	}
}

func newDomainCommand(msg *tgbotapi.Message) models.UserCommand {
	return models.UserCommand{
		Command:     models.Command(msg.Command()),
		UserMessage: newDomainMessage(msg),
	}
}
