package storage

import (
	"context"

	"github.com/gotd/td/telegram/uploader"
	"go.uber.org/zap"
)

var _ uploader.Progress = (*uploadProgress)(nil)

type uploadProgress struct {
	logger  *zap.Logger
	percent float32
}

func newUploadProgress(logger *zap.Logger) *uploadProgress {
	return &uploadProgress{logger: logger}
}

func (u *uploadProgress) Chunk(ctx context.Context, state uploader.ProgressState) error {
	u.percent = float32(state.Uploaded) / float32(state.Total)
	u.logger.Info("upload progress update", zap.Float32("percentage", u.percent))

	return nil
}
