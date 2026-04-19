package curl

import (
	"context"
	"fmt"
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
	if len(args) == 0 {
		return fmt.Errorf("no URL specified")
	}

	// Simple simulation of curl output
	fmt.Fprintf(c.stdout, "Simulating curl fetch for: %s\n", args[0])
	return nil
}
