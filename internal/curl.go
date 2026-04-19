package internal

import (
	"context"
	"io"
)

type Command struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func New(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Command {
	return &Command{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *Command) Execute(ctx context.Context, args []string) error {
	return nil
}
