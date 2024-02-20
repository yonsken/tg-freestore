package storage

import (
	"context"

	"github.com/gotd/td/telegram/uploader"
)

var _ uploader.Progress = (*uploadProgress)(nil)

type uploadProgress struct {
	percent float32
}

func (u *uploadProgress) Chunk(ctx context.Context, state uploader.ProgressState) error {
	u.percent = float32(state.Uploaded) / float32(state.Total)
	return nil
}
