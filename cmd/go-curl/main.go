package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yarencheng/go-curl/internal"
)

func main() {
	// Configure zerolog for pretty console output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	ctx := context.Background()
	args := os.Args[1:]

	cmd := internal.New(os.Stdin, os.Stdout, os.Stderr, log.Logger)
	if err := cmd.Execute(ctx, args); err != nil {
		log.Error().Err(err).Msg("Command failed")
		os.Exit(1)
	}
}
