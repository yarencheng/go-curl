package curl_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/yarencheng/go-curl/pkg/curl"
)

func TestNew(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := io.Discard
	stderr := io.Discard
	logger := zerolog.Nop()
	fs := afero.NewMemMapFs()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
	assert.NotNil(t, cmd)
}

func TestExecute(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := io.Discard
	stderr := io.Discard
	logger := zerolog.Nop()
	fs := afero.NewMemMapFs()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
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
	fs := afero.NewMemMapFs()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
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
	fs := afero.NewMemMapFs()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
	ctx := context.Background()
	args := []string{"-s", "-v", "-i", "https://example.com"}

	err := cmd.Execute(ctx, args)
	assert.NoError(t, err)
}

func TestExecute_IPFlags(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()
	fs := afero.NewMemMapFs()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
	ctx := context.Background()

	// Test -4
	args4 := []string{"-4", "https://example.com"}
	err := cmd.Execute(ctx, args4)
	assert.NoError(t, err)

	// Test -6
	args6 := []string{"-6", "https://example.com"}
	err = cmd.Execute(ctx, args6)
	assert.NoError(t, err)
}

func TestExecute_Complex(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()
	fs := afero.NewMemMapFs()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
	ctx := context.Background()
	args := []string{"-X", "POST", "-H", "Content-Type: application/json", "-d", `{"key":"value"}`, "https://httpbin.org/post"}

	err := cmd.Execute(ctx, args)
	assert.NoError(t, err)
}

func TestExecute_FileFlags(t *testing.T) {
	headerFile := "test_headers.txt"
	dataFile := "test_data.json"
	
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, headerFile, []byte("X-Test: test-value\n"), 0644)
	assert.NoError(t, err)
	
	err = afero.WriteFile(fs, dataFile, []byte(`{"hello":"world"}`), 0644)
	assert.NoError(t, err)

	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
	ctx := context.Background()
	args := []string{"-X", "POST", "-H", "@" + headerFile, "-d", "@" + dataFile, "https://httpbin.org/post"}

	err = cmd.Execute(ctx, args)
	assert.NoError(t, err)
}

func TestExecute_CookiesAndUpload(t *testing.T) {
	cookieJar := "test_cookies.txt"
	uploadFile := "test_upload.txt"
	
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, uploadFile, []byte("upload test data"), 0644)
	assert.NoError(t, err)

	stdin := strings.NewReader("")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := zerolog.Nop()

	cmd := curl.New(stdin, stdout, stderr, logger, fs)
	ctx := context.Background()
	
	// Test upload
	args := []string{"-T", uploadFile, "https://httpbin.org/put"}
	err = cmd.Execute(ctx, args)
	assert.NoError(t, err)
	
	// Test cookies
	args = []string{"-b", "name=value", "-c", cookieJar, "https://httpbin.org/cookies/set/test/val"}
	err = cmd.Execute(ctx, args)
	assert.NoError(t, err)
	
	content, err := afero.ReadFile(fs, cookieJar)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test")
	assert.Contains(t, string(content), "val")
}
