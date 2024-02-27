package storage

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/yonsken/tg-freestore/internal/client"
	"go.uber.org/zap"
)

type Manager struct {
	logger     *zap.Logger
	client     *telegram.Client
	stopClient bg.StopFunc
}

func NewManager(logger *zap.Logger) (*Manager, error) {
	client, stopClient, err := client.NewClient(logger)
	if err != nil {
		return nil, fmt.Errorf("creating new client: %w", err)
	}

	manager := Manager{
		logger:     logger,
		client:     client,
		stopClient: stopClient,
	}

	return &manager, nil
}

func (m *Manager) UploadFile(ctx context.Context, filePath string, recipient string) error {
	m.logger.Info("uploading file", zap.String("filePath", filePath))

	api := m.client.API()

	uploader := uploader.NewUploader(api).WithProgress(newUploadProgress(m.logger))

	upload, err := uploader.FromPath(ctx, filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("uploading file: %w", err)
	}

	file := message.File(upload).Filename(filepath.Base(filePath))

	if _, err := message.NewSender(api).Resolve(recipient).Media(ctx, file); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (m *Manager) Close() {
	if err := m.stopClient(); err != nil {
		m.logger.Error("failed to close storage manager", zap.Error(err))
	}
}
