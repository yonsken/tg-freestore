package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yonsken/tg-freestore/internal/env"
	"github.com/yonsken/tg-freestore/internal/storage"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 1 {
		return nil
	}

	vars, err := env.GetRequiredEnvVars()
	if err != nil {
		return fmt.Errorf("getting required environment variables: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	manager, err := storage.NewManager(logger)
	if err != nil {
		return fmt.Errorf("creating storage manager: %w", err)
	}
	defer manager.Close()

	ctx := context.Background()

	switch os.Args[0] {
	case "list":
		if err := manager.ListFiles(ctx, vars.UploadFileRecipient); err != nil {
			return fmt.Errorf("listing test files: %w", err)
		}
	case "upload":
		if err := manager.UploadFile(ctx, vars.UploadFilePath, vars.UploadFileRecipient); err != nil {
			return fmt.Errorf("uploading test file: %w", err)
		}
	}

	return nil
}
