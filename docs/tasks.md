# Functional Parity Tracking: curl

This document tracks the alignment of the `go-curl` implementation with upstream `curl`.

## Overview

Status codes:
- [x] : Fully implemented and verified.
- [ ] : Missing or incomplete.
- [-] : Deliberately skipped.

---

## curl Flags

### Basic Flags
- [x] `-X, --request <command>`: Specify request command to use (upstream: `third_party/curl/src/tool_getparam.c:2796`)
- [x] `-d, --data <data>`: HTTP POST data (upstream: `third_party/curl/src/tool_getparam.c:2095`)
- [x] `-H, --header <header/@file>`: Pass custom header(s) to server (upstream: `third_party/curl/src/tool_getparam.c:2438`)
- [x] `-i, --include`: Include protocol response headers in the output (upstream: `third_party/curl/src/tool_getparam.c:2479`)
- [x] `-u, --user <user:password>`: Server user and password (upstream: `third_party/curl/src/tool_getparam.c:2722`)
- [x] `-o, --output <file>`: Write to file instead of stdout (upstream: `third_party/curl/src/tool_getparam.c:2558`)
- [x] `-O, --remote-name`: Write output to a file named as the remote file (upstream: `third_party/curl/src/tool_getparam.c:2559`)
- [x] `-v, --verbose`: Make the operation more talkative (upstream: `third_party/curl/src/tool_getparam.c:2732`)
- [x] `-s, --silent`: Silent mode (upstream: `third_party/curl/src/tool_getparam.c:2679`)
- [x] `-L, --location`: Follow redirects (upstream: `third_party/curl/src/tool_getparam.c:2517`)
- [x] `--version`: Show version number and exit (upstream: `third_party/curl/src/tool_getparam.c:2750`)
- [-] Flag `--interface`: Not sandbox-friendly (requires binding to specific IP) (upstream: `third_party/curl/src/tool_getparam.c:1592`)
- [-] Flag `--dns-servers`: Browser/WASM sandbox limitation (upstream: `third_party/curl/src/tool_getparam.c:1469`)
- [-] Flag `--unix-socket`: OS-level socket control unavailable (upstream: `third_party/curl/src/tool_getparam.c:1907`)
- [-] Flag `--proxy`: Handled by browser environment, not programmatically via fetch (upstream: `third_party/curl/src/tool_getparam.c:2790`)

### Advanced Flags
- [x] `-A, --user-agent <name>`: Send User-Agent <name> to server (upstream: `third_party/curl/src/tool_getparam.c:2051`)
- [x] `-b, --cookie <data>`: Send cookies from string/file (upstream: `third_party/curl/src/tool_getparam.c:2066`)
- [x] `-c, --cookie-jar <filename>`: Write cookies to <filename> after operation (upstream: `third_party/curl/src/tool_getparam.c:2080`)
- [x] `-e, --referer <url>`: Referrer URL (upstream: `third_party/curl/src/tool_getparam.c:2109`)
- [x] `-f, --fail`: Fail silently (no output at all) on HTTP errors (upstream: `third_party/curl/src/tool_getparam.c:2395`)
- [x] `-m, --max-time <seconds>`: Maximum time allowed for the transfer (upstream: `third_party/curl/src/tool_getparam.c:2520`)
- [x] `-T, --upload-file <file>`: Transfer local FILE to destination (upstream: `third_party/curl/src/tool_getparam.c:2689`)
