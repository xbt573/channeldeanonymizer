package bot

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type Bot struct {
	options Options
}

type Options struct {
	Token  string
	ChatId int64
}

func New(options Options) *Bot {
	return &Bot{options}
}

func (b *Bot) onMessage(bot *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Id != b.options.ChatId {
		return nil
	}

	administrators, err := ctx.EffectiveChat.GetAdministrators(bot, nil)
	if err != nil {
		return err
	}

	names := []string{}
	for _, administrator := range administrators {
		member := administrator.GetUser()

		name := strings.TrimSpace(fmt.Sprintf("%v %v", member.FirstName, member.LastName))
		names = append(names, name)
	}

	if slices.Contains(names, ctx.EffectiveSender.AuthorSignature) {
		return nil
	}

	_, err = ctx.EffectiveMessage.Delete(bot, nil)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) Start() error {
	instance, err := gotgbot.NewBot(b.options.Token, nil)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},

		MaxRoutines: ext.DefaultMaxRoutines,
		Logger:      logger,
	})

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{Logger: logger})

	dispatcher.AddHandler(handlers.NewMessage(message.All, b.onMessage).SetAllowChannel(true))

	err = updater.StartPolling(instance, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		return err
	}

	updater.Idle()
	return nil // probably unreachable
}
