package curl_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/yarencheng/go-curl/pkg/curl"
)

func TestExecute_Errors(t *testing.T) {
	stdin := strings.NewReader("")
	logger := zerolog.Nop()
	fs := afero.NewMemMapFs()
	cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)

	t.Run("no URL", func(t *testing.T) {
		err := cmd.Execute(context.Background(), []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no URL specified")
	})

	t.Run("invalid flag", func(t *testing.T) {
		err := cmd.Execute(context.Background(), []string{"--invalid-flag", "http://localhost"})
		assert.Error(t, err)
	})

	t.Run("missing data file", func(t *testing.T) {
		err := cmd.Execute(context.Background(), []string{"-d", "@nonexistent.txt", "http://localhost"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open data file")
	})

	t.Run("missing header file", func(t *testing.T) {
		err := cmd.Execute(context.Background(), []string{"-H", "@nonexistent.txt", "http://localhost"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read header file")
	})
    
    t.Run("missing upload file", func(t *testing.T) {
		err := cmd.Execute(context.Background(), []string{"-T", "nonexistent.txt", "http://localhost"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open upload file")
	})
    
}

func TestExecute_MockServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Path == "/cookies" {
			http.SetCookie(w, &http.Cookie{Name: "testname", Value: "testvalue", Path: "/"})
		}
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			if string(body) == "request body content" {
				w.Header().Set("X-Received-Body", "true")
			}
		}
		w.Header().Set("X-Custom-Resp", "val")
		fmt.Fprint(w, "response body")
	}))
	defer ts.Close()

	stdin := strings.NewReader("")
	logger := zerolog.Nop()
	fs := afero.NewMemMapFs()

	t.Run("include headers", func(t *testing.T) {
		stdout := &strings.Builder{}
		cmd := curl.New(stdin, stdout, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-i", ts.URL})
		assert.NoError(t, err)
		assert.Contains(t, stdout.String(), "X-Custom-Resp: val")
		assert.Contains(t, stdout.String(), "response body")
	})

	t.Run("verbose mode", func(t *testing.T) {
		stderr := &strings.Builder{}
		cmd := curl.New(stdin, io.Discard, stderr, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-v", ts.URL})
		assert.NoError(t, err)
		assert.Contains(t, stderr.String(), "* Connected to")
		assert.Contains(t, stderr.String(), "> GET")
		assert.Contains(t, stderr.String(), "< HTTP/1.1 200 OK")
	})

	t.Run("fail flag", func(t *testing.T) {
		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-f", ts.URL + "/error"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error: 404")
	})

	t.Run("post data from file", func(t *testing.T) {
		dataFile := "test_data.txt"
		afero.WriteFile(fs, dataFile, []byte("request body content"), 0644)

		stdout := &strings.Builder{}
		cmd := curl.New(stdin, stdout, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-d", "@" + dataFile, "-i", ts.URL})
		assert.NoError(t, err)
		assert.Contains(t, stdout.String(), "X-Received-Body: true")
	})

	t.Run("custom user agent and referer", func(t *testing.T) {
		// Just verify it doesn't crash, we'd need to mock the server to check headers
		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-A", "my-ua", "-e", "http://ref.com", ts.URL})
		assert.NoError(t, err)
	})

	t.Run("basic auth", func(t *testing.T) {
		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-u", "user:pass", ts.URL})
		assert.NoError(t, err)
		
		err = cmd.Execute(context.Background(), []string{"-u", "useronly", ts.URL})
		assert.NoError(t, err)
	})

	t.Run("cookie handling", func(t *testing.T) {
		cookieJar := "test_jar.txt"
		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		
		// Set cookie
		err := cmd.Execute(context.Background(), []string{"-c", cookieJar, ts.URL + "/cookies"})
		assert.NoError(t, err)
		
		// Use cookie from jar
		err = cmd.Execute(context.Background(), []string{"-b", "name=value", ts.URL})
		assert.NoError(t, err)
	})

	t.Run("netscape cookie file", func(t *testing.T) {
		cookieFile := "test_netscape_cookies.txt"
		// domain, flag, path, secure, expiration, name, value
		content := "localhost\tTRUE\t/\tFALSE\t0\tmysession\tabcdef\n"
		afero.WriteFile(fs, cookieFile, []byte(content), 0644)

		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-b", cookieFile, ts.URL})
		assert.NoError(t, err)
	})

	t.Run("cookie string", func(t *testing.T) {
		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-b", "name1=val1;name2=val2", ts.URL})
		assert.NoError(t, err)
	})

	t.Run("remote name failed", func(t *testing.T) {
		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-O", ts.URL + "/"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not determine remote name")
	})
    
    t.Run("header from file multi-line", func(t *testing.T) {
        headerFile := "test_headers_multi.txt"
		content := "X-Header-1: val1\nX-Header-2: val2\n\n"
		afero.WriteFile(fs, headerFile, []byte(content), 0644)

		cmd := curl.New(stdin, io.Discard, io.Discard, logger, fs)
		err := cmd.Execute(context.Background(), []string{"-H", "@" + headerFile, ts.URL})
		assert.NoError(t, err)
    })
}
