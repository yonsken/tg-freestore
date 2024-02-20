package storage

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

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

	progress := new(uploadProgress)
	progressCtx, progressCancel := context.WithCancel(ctx)
	defer progressCancel()

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			select {
			case <-progressCtx.Done():
				return
			case <-ticker.C:
				m.logger.Info("upload progress update", zap.Float32("percentage", progress.percent))
			}
		}
	}()

	uploader := uploader.NewUploader(api).WithProgress(progress)

	upload, err := uploader.FromPath(ctx, filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("uploading file: %w", err)
	}

	progressCancel()

	file := message.File(upload).Filename(filepath.Base(filePath))

	if _, err := message.NewSender(api).Resolve(recipient).Media(ctx, file); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (m *Manager) Close() error {
	return m.stopClient()
}
