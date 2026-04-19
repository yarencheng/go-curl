package internal_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/yarencheng/go-curl/internal"
)

func TestNew(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := io.Discard
	stderr := io.Discard
	logger := zerolog.Nop()

	cmd := internal.New(stdin, stdout, stderr, logger)
	assert.NotNil(t, cmd)
}

func TestExecute(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := io.Discard
	stderr := io.Discard
	logger := zerolog.Nop()

	cmd := internal.New(stdin, stdout, stderr, logger)
	ctx := context.Background()
	args := []string{"https://example.com"}

	err := cmd.Execute(ctx, args)
	assert.NoError(t, err)
}
