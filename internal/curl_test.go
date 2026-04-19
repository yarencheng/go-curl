package internal_test

import (
	"context"
	"io"
	"os"
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
func TestExecute_Version(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()

	cmd := internal.New(stdin, stdout, stderr, logger)
	ctx := context.Background()
	args := []string{"--version"}

	err := cmd.Execute(ctx, args)
	assert.NoError(t, err)
	assert.Contains(t, stdout.String(), "go-curl version")
}

func TestExecute_Flags(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()

	cmd := internal.New(stdin, stdout, stderr, logger)
	ctx := context.Background()
	args := []string{"-s", "-v", "-i", "https://example.com"}

	err := cmd.Execute(ctx, args)
	assert.NoError(t, err)
}

func TestExecute_Complex(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()

	cmd := internal.New(stdin, stdout, stderr, logger)
	ctx := context.Background()
	args := []string{"-X", "POST", "-H", "Content-Type: application/json", "-d", `{"key":"value"}`, "https://httpbin.org/post"}

	err := cmd.Execute(ctx, args)
	assert.NoError(t, err)
}

func TestExecute_FileFlags(t *testing.T) {
	headerFile := "test_headers.txt"
	dataFile := "test_data.json"
	
	err := os.WriteFile(headerFile, []byte("X-Test: test-value\n"), 0644)
	assert.NoError(t, err)
	defer os.Remove(headerFile)
	
	err = os.WriteFile(dataFile, []byte(`{"hello":"world"}`), 0644)
	assert.NoError(t, err)
	defer os.Remove(dataFile)

	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()

	cmd := internal.New(stdin, stdout, stderr, logger)
	ctx := context.Background()
	args := []string{"-X", "POST", "-H", "@" + headerFile, "-d", "@" + dataFile, "https://httpbin.org/post"}

	err = cmd.Execute(ctx, args)
	assert.NoError(t, err)
}

func TestExecute_CookiesAndUpload(t *testing.T) {
	cookieJar := "test_cookies.txt"
	uploadFile := "test_upload.txt"
	
	err := os.WriteFile(uploadFile, []byte("upload test data"), 0644)
	assert.NoError(t, err)
	defer os.Remove(uploadFile)
	defer os.Remove(cookieJar)

	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()

	cmd := internal.New(stdin, stdout, stderr, logger)
	ctx := context.Background()
	
	// Test upload
	args := []string{"-T", uploadFile, "https://httpbin.org/put"}
	err = cmd.Execute(ctx, args)
	assert.NoError(t, err)
	
	// Test cookies
	args = []string{"-b", "name=value", "-c", cookieJar, "https://httpbin.org/cookies/set/test/val"}
	err = cmd.Execute(ctx, args)
	assert.NoError(t, err)
	
	content, err := os.ReadFile(cookieJar)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test")
	assert.Contains(t, string(content), "val")
}
