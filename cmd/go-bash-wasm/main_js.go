//go:build js

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"syscall/js"

	"github.com/rs/zerolog"
	"github.com/yarencheng/go-curl/internal/commands/curl"
)

func main() {
	// zerolog is redirected to browser console via BrowserConsoleWriter.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	ctx := context.Background()

	// Create a synchronous pipe. The JS side writes lines in; the Go
	// logic reads them out.
	pr, pw := io.Pipe()

	// Expose writeStdin(line string) to JS.
	writeStdinFn := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		data := args[0].String()
		go func() {
			if _, err := pw.Write([]byte(data)); err != nil {
				fmt.Fprintf(os.Stderr, "writeStdin: %v\n", err)
			}
		}()
		return nil
	})
	js.Global().Set("writeStdin", writeStdinFn)
	defer writeStdinFn.Release()

	// Expose closeStdin() to JS.
	closeStdinFn := js.FuncOf(func(_ js.Value, _ []js.Value) any {
		pw.Close()
		return nil
	})
	js.Global().Set("closeStdin", closeStdinFn)
	defer closeStdinFn.Release()

	c := curl.New(pr, os.Stdout, os.Stderr)
	// In WASM, we might want to pass args from JS global or similar.
	// For now, mirroring go-bash-wasm pattern.
	if err := c.Execute(ctx, nil); err != nil {
		fmt.Fprintf(os.Stderr, "curl simulator failed: %v\n", err)
	}

	// Keep the WASM module alive after exit so the page
	// doesn't throw "Go program has already exited" on callbacks.
	select {}
}
