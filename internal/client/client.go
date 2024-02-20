package client

import (
	"fmt"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func NewClient(logger *zap.Logger) (*telegram.Client, bg.StopFunc, error) {
	dispatcher := tg.NewUpdateDispatcher()

	opts := telegram.Options{
		Logger:        logger,
		UpdateHandler: dispatcher,
	}

	client, err := telegram.ClientFromEnvironment(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("creating client from environment: %w", err)
	}

	stop, err := bg.Connect(client)
	if err != nil {
		return nil, nil, fmt.Errorf("connecting client: %w", err)
	}

	return client, stop, nil
}
