package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yarencheng/go-curl/pkg/curl"
)

func main() {
	// Configure zerolog for pretty console output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	ctx := context.Background()
	args := os.Args[1:]

	cmd := curl.New(os.Stdin, os.Stdout, os.Stderr, log.Logger, curl.DefaultFs())
	if err := cmd.Execute(ctx, args); err != nil {
		log.Error().Err(err).Msg("Command failed")
		os.Exit(1)
	}
}
