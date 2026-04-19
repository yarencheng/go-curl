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
- [ ] `-X, --request <command>`: Specify request command to use
- [ ] `-d, --data <data>`: HTTP POST data
- [ ] `-H, --header <header/@file>`: Pass custom header(s) to server
- [ ] `-i, --include`: Include protocol response headers in the output
- [ ] `-u, --user <user:password>`: Server user and password
- [ ] `-o, --output <file>`: Write to file instead of stdout
- [ ] `-O, --remote-name`: Write output to a file named as the remote file
- [ ] `-v, --verbose`: Make the operation more talkative
- [ ] `-s, --silent`: Silent mode
- [ ] `-L, --location`: Follow redirects
- [ ] `--version`: Show version number and exit

### Advanced Flags
- [ ] `-A, --user-agent <name>`: Send User-Agent <name> to server
- [ ] `-b, --cookie <data>`: Send cookies from string/file
- [ ] `-c, --cookie-jar <filename>`: Write cookies to <filename> after operation
- [ ] `-e, --referer <url>`: Referrer URL
- [ ] `-f, --fail`: Fail silently (no output at all) on HTTP errors
- [ ] `-m, --max-time <seconds>`: Maximum time allowed for the transfer
- [ ] `-T, --upload-file <file>`: Transfer local FILE to destination
