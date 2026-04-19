package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

const version = "0.1.0"

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
func (c *Command) Execute(ctx context.Context, args []string) error {
	fs := pflag.NewFlagSet("curl", pflag.ContinueOnError)
	fs.SetOutput(c.stderr)

	showVersion := fs.Bool("version", false, "Show version number and exit")
	silent := fs.BoolP("silent", "s", false, "Silent mode")
	verbose := fs.BoolP("verbose", "v", false, "Make the operation more talkative")
	include := fs.BoolP("include", "i", false, "Include protocol response headers in the output")
	request := fs.StringP("request", "X", "", "Specify request command to use")
	data := fs.StringP("data", "d", "", "HTTP POST data")
	headers := fs.StringArrayP("header", "H", []string{}, "Pass custom header(s) to server")
	output := fs.StringP("output", "o", "", "Write to file instead of stdout")
	remoteName := fs.BoolP("remote-name", "O", false, "Write output to a file named as the remote file")
	user := fs.StringP("user", "u", "", "Server user and password")
	location := fs.BoolP("location", "L", false, "Follow redirects")
	userAgent := fs.StringP("user-agent", "A", "", "Send User-Agent <name> to server")
	referer := fs.StringP("referer", "e", "", "Referer URL")
	fail := fs.BoolP("fail", "f", false, "Fail silently (no output at all) on HTTP errors")
	maxTime := fs.Float64P("max-time", "m", 0, "Maximum time allowed for the transfer")
	cookie := fs.StringP("cookie", "b", "", "Send cookies from string/file")
	cookieJar := fs.StringP("cookie-jar", "c", "", "Write cookies to <filename> after operation")
	uploadFile := fs.StringP("upload-file", "T", "", "Transfer local FILE to destination")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *showVersion {
		fmt.Fprintf(c.stdout, "go-curl version %s\n", version)
		return nil
	}

	remainingArgs := fs.Args()
	if len(remainingArgs) == 0 {
		return fmt.Errorf("no URL specified")
	}
	url := remainingArgs[0]

	method := *request
	if method == "" {
		if *data != "" {
			method = http.MethodPost
		} else {
			method = http.MethodGet
		}
	}

	var body io.Reader
	if *data != "" {
		if strings.HasPrefix(*data, "@") {
			fileName := (*data)[1:]
			f, err := os.Open(fileName)
			if err != nil {
				return fmt.Errorf("failed to open data file %s: %w", fileName, err)
			}
			defer f.Close()
			body = f
		} else {
			body = strings.NewReader(*data)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for _, h := range *headers {
		if strings.HasPrefix(h, "@") {
			fileName := h[1:]
			content, err := os.ReadFile(fileName)
			if err != nil {
				return fmt.Errorf("failed to read header file %s: %w", fileName, err)
			}
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
		} else {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	if *userAgent != "" {
		req.Header.Set("User-Agent", *userAgent)
	}
	if *referer != "" {
		req.Header.Set("Referer", *referer)
	}

	if *user != "" {
		parts := strings.SplitN(*user, ":", 2)
		if len(parts) == 2 {
			req.SetBasicAuth(parts[0], parts[1])
		} else {
			req.SetBasicAuth(parts[0], "")
		}
	}

	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	if *cookie != "" {
		if strings.Contains(*cookie, "=") && !strings.Contains(*cookie, "/") {
			// Simple name=value string
			parts := strings.Split(*cookie, ";")
			var cookies []*http.Cookie
			for _, p := range parts {
				kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
				if len(kv) == 2 {
					cookies = append(cookies, &http.Cookie{Name: kv[0], Value: kv[1]})
				}
			}
			u, _ := req.URL.Parse(url)
			client.Jar.SetCookies(u, cookies)
		} else {
			// filename
			content, err := os.ReadFile(*cookie)
			if err == nil {
				// Simple parsing of cookie file (Netscape format is complex, doing simple for now)
				lines := strings.Split(string(content), "\n")
				var cookies []*http.Cookie
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "#") {
						continue
					}
					parts := strings.Fields(line)
					if len(parts) >= 7 {
						cookies = append(cookies, &http.Cookie{
							Name:  parts[5],
							Value: parts[6],
						})
					}
				}
				u, _ := req.URL.Parse(url)
				client.Jar.SetCookies(u, cookies)
			}
		}
	}

	if *uploadFile != "" {
		method = http.MethodPut
		f, err := os.Open(*uploadFile)
		if err != nil {
			return fmt.Errorf("failed to open upload file %s: %w", *uploadFile, err)
		}
		defer f.Close()
		req.Method = method
		req.Body = io.NopCloser(f)
	}

	if *maxTime > 0 {
		client.Timeout = time.Duration(*maxTime * float64(time.Second))
	}

	if !*location {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if *fail && resp.StatusCode >= 400 {
		return fmt.Errorf("The requested URL returned error: %d", resp.StatusCode)
	}

	if *verbose {
		fmt.Fprintf(c.stderr, "* Connected to %s\n", url)
		fmt.Fprintf(c.stderr, "> %s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
		for name, values := range req.Header {
			for _, value := range values {
				fmt.Fprintf(c.stderr, "> %s: %s\n", name, value)
			}
		}
		fmt.Fprintln(c.stderr, ">")
	}

	var out io.Writer = c.stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		out = f
	} else if *remoteName {
		fileName := path.Base(req.URL.Path)
		if fileName == "/" || fileName == "." || fileName == "" {
			return fmt.Errorf("could not determine remote name from URL")
		}
		f, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	if *include {
		fmt.Fprintf(out, "%s %s\n", resp.Proto, resp.Status)
		for name, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(out, "%s: %s\n", name, value)
			}
		}
		fmt.Fprintln(out)
	}

	if !*silent {
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
	} else {
		_, _ = io.Copy(io.Discard, resp.Body)
	}

	if *cookieJar != "" {
		u, _ := req.URL.Parse(url)
		cookies := client.Jar.Cookies(u)
		var sb strings.Builder
		sb.WriteString("# Netscape HTTP Cookie File\n")
		for _, c := range cookies {
			// domain, flag, path, secure, expiration, name, value
			fmt.Fprintf(&sb, "%s\tTRUE\t%s\tFALSE\t0\t%s\t%s\n", req.URL.Host, c.Path, c.Name, c.Value)
		}
		err := os.WriteFile(*cookieJar, []byte(sb.String()), 0644)
		if err != nil {
			return fmt.Errorf("failed to write cookie jar: %w", err)
		}
	}

	return nil
}
