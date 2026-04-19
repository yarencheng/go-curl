package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/yarencheng/go-curl/internal/commands/curl"
)

func main() {
	// Initialize logging
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	ctx := logger.WithContext(context.Background())

	c := curl.New(os.Stdin, os.Stdout, os.Stderr)
	if err := c.Execute(ctx, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "curl failed: %v\n", err)
		os.Exit(1)
	}
}
