package bot

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// Required environment variables:
//
//	BOT_TOKEN:     token from BotFather
//	APP_ID:        app_id of Telegram app
//	APP_HASH:      app_hash of Telegram app
func RunBot() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("building logger: %w", err)
	}

	dispatcher := tg.NewUpdateDispatcher()

	opts := telegram.Options{
		Logger:        logger,
		UpdateHandler: dispatcher,
	}

	if err := telegram.BotFromEnvironment(
		context.Background(),
		opts,
		setupBot(dispatcher),
		telegram.RunUntilCanceled,
	); err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	return nil
}

type setupBotFunc func(ctx context.Context, client *telegram.Client) error

func setupBot(dispatcher tg.UpdateDispatcher) setupBotFunc {
	return func(ctx context.Context, client *telegram.Client) error {
		var (
			api    = tg.NewClient(client)
			sender = message.NewSender(api)
		)

		onNewMessageFunc := func(
			ctx context.Context,
			entities tg.Entities,
			update *tg.UpdateNewMessage,
		) error {
			m, ok := update.Message.(*tg.Message)
			if !ok || m.Out {
				return nil
			}

			_, err := sender.Reply(entities, update).Text(ctx, m.Message)
			return err
		}

		dispatcher.OnNewMessage(onNewMessageFunc)
		return nil
	}
}
