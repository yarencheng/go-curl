package internal

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

// Command represents the curl command execution environment.
type Command struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	logger zerolog.Logger
}

// New creates a new curl Command with the provided I/O streams and default logger.
func New(stdin io.Reader, stdout io.Writer, stderr io.Writer, logger zerolog.Logger) *Command {
	return &Command{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		logger: logger,
	}
}

// Execute runs the curl command with the provided arguments.
// Currently, it logs the arguments and returns nil.
func (c *Command) Execute(ctx context.Context, args []string) error {
	c.logger.Info().
		Strs("args", args).
		Msg("Executing curl command")

	return nil
}
