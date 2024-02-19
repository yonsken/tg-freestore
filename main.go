package main

import (
	"fmt"
	"os"

	"github.com/yonsken/tg-freestore/internal/bot"
)

func main() {
	if err := bot.RunBot(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
