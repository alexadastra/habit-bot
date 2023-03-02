package internal

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type UserMessage struct {
	Id      int64
	Message string
}

type UserCommand struct {
	Command string
	UserMessage
}

type Bot struct {
	botAPI *tgbotapi.BotAPI
	inCh   tgbotapi.UpdatesChannel

	cancel context.CancelFunc
	wg     *sync.WaitGroup

	commandCh chan UserCommand
	messageCh chan UserMessage
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
		commandCh: make(chan UserCommand),
		messageCh: make(chan UserMessage),
	}, nil
}

func (b *Bot) GetCommandsChan() chan UserCommand {
	return b.commandCh
}

func (b *Bot) GetMessagesChan() chan UserMessage {
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
			// Handle commands
			if update.Message != nil && update.Message.IsCommand() {
				b.commandCh <- UserCommand{
					Command:     update.Message.Command(),
					UserMessage: UserMessage{Id: update.Message.From.ID, Message: update.Message.Text},
				}
			}

			// Handle messages that are not commands
			if update.Message != nil && !update.Message.IsCommand() {
				b.messageCh <- UserMessage{Id: update.Message.From.ID, Message: update.Message.Text}
			}
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
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.botAPI.Send(msg)
	return err
}
