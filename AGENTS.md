# AGENTS.md

Guide for AI agents (and humans) working on Candy QR.

## What this is

Candy QR is a Linux terminal (TUI) app for creating QR codes. It covers
three encodings: Wi-Fi credentials, URLs, and vCards. It can print the QR to
the terminal (when the terminal supports it) and export to PNG.

## Stack

- **Language:** Go
- **TUI framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
  (Elm Architecture: `Model`, `Init`, `Update`, `View`)
- **Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss)
  (added later, in a design pass)
- **Components:** [Bubbles](https://github.com/charmbracelet/bubbles)
  (text inputs, spinners, etc.)
- **QR generation:** `github.com/skip2/go-qrcode` (or equivalent, see below)
- **PNG export:** standard `image` + `image/png`

## Architecture

The app follows the Bubble Tea Elm Architecture:

```
Model    — app state (current mode, form fields)
Init     — initial command (returns nil unless I/O is needed)
Update   — handle Msg (keypresses, form submissions) -> (Model, Cmd)
View     — render the UI as a string
```

Keep state in the `Model`. Treat every interaction as a message. Prefer
declarative views over manual screen redraws.

## Product focus

The **vCard** flow is the primary feature: the goal is to let someone hold up
their screen and have another person's camera app drop the contact into their
address book. Design around that happy path first.

The three modes are:

1. `url` — a plain URL encoded as a QR code
2. `wifi` — `WIFI:S:<ssid>;T:WPA;P:<password>;;`
3. `vcard` — `BEGIN:VCARD ... END:VCARD`

## Output

Two output paths are required:

- **Terminal print** — render the QR as Unicode blocks to stdout/TUI,
  detecting terminal capability.
- **PNG export** — save the QR to a PNG file.

## Conventions

- Follow the existing Go style (`gofmt`, no em dash in source code).
- Match Charm's brand voice in docs and copy: playful, warm, benefit-driven.
  Keep the whimsy in docs/UI, not in logic.
- Do not add comments unless they explain _why_.
- Do not commit unless explicitly asked.

## Testing

- Run `go test ./...` after changes.
- Run `go vet ./...` and `gofmt -l .` when available.

## Scope guardrails

- This is **v0.1**: functionality first, no colors or special design. Colors,
  gradients, rounded corners, and logos are explicitly out of scope until a
  later design pass.
- Keep the QR payload encoding and the rendering as separate concerns so the
  "pretty" pass can be layered on without rewriting the core.