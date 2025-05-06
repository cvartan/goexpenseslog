package bot

import (
	"context"
	"fmt"
	"log"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotRequest struct {
	botapi.Update
}

type BotResponse struct {
	botapi.MessageConfig
}

type BotHandleFunc func(*BotRequest, *BotResponse) error

type TelegramBot struct {
	token           string
	commandHandlers map[string]BotHandleFunc
	defaultHandler  BotHandleFunc
}

func New(token string) *TelegramBot {
	if token == "" {
		panic("token must be defined")
	}
	return &TelegramBot{
		token:           token,
		commandHandlers: make(map[string]BotHandleFunc, 8),
	}
}

func (bot *TelegramBot) SetDefaultHandler(handler BotHandleFunc) {
	if handler == nil {
		panic("default handler must be defined")
	}
	bot.defaultHandler = handler
}

func (bot *TelegramBot) SetCommandHandler(command string, handler BotHandleFunc) {
	if command == "" {
		panic("handled command must be defined")
	}
	if handler == nil {
		panic("handelr for command is not be nil")
	}
	bot.commandHandlers[command] = handler
}

func (bot *TelegramBot) Listen(ctx context.Context) error {
	b, err := botapi.NewBotAPI(bot.token)
	if err != nil {
		return fmt.Errorf("creating bot error: %v", err)
	}

	updater := botapi.NewUpdate(0)
	updater.Timeout = 60

	messageUpdates := b.GetUpdatesChan(updater)

	messages := make(chan botapi.Update)

	go bot.listen(messages, messageUpdates)

	for {
		select {
		case msg := <-messages:
			{
				cmd := msg.Message.Command()
				var handler BotHandleFunc
				handler, ok := bot.commandHandlers[cmd]
				if !ok {
					handler = bot.defaultHandler
				}
				if handler != nil {
					go func() {
						req := &BotRequest{msg}
						resp := &BotResponse{botapi.NewMessage(msg.Message.Chat.ID, "")}
						if err := handler(req, resp); err != nil {
							log.Println(err)
						}
						resp.ParseMode = botapi.ModeHTML
						b.Send(resp.MessageConfig)
					}()
				}
			}
		case <-ctx.Done():
			{
				close(messages)
				return nil
			}
		}
	}
}

func (bot *TelegramBot) listen(messages chan<- botapi.Update, updates botapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			messages <- update
		}
	}
}
