# 🍬 Candy QR

> Your terminal, but make it a QR code business card.

Candy QR is a Linux terminal app for making beautiful, scannable QR codes.
Build a Wi-Fi login, a link, or a vCard, then print it straight to your
terminal or export it as a PNG.

Open your camera, point it at the screen, and hand someone your contact info
with _pizzazz_ 💅🏻 — no typing, no fumbling, just a scan.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), the fun,
functional and stateful way to build terminal apps.

## Features

- **vCard** — encode your name, phone, email, and more directly into a
  scannable contact that drops into a phone's address book
- **Wi-Fi** — share your network name and password as a scannable code
- **URL** — turn any link into a QR code
- **Live preview** — watch the QR update in real time as you type
- **Export to PNG** — save a crisp, shareable image

## Status

**v0.1** — the delicious-but-plain first cut. Functionality first, no colors
or fancy design _yet_. Logo, gradients, rounded corners, and the full glam
treatment are on the roadmap.

## Install

```bash
go install github.com/RM35M4419/candy-qr@latest
```

## Usage

One command, one flow:

```
candy-qr
```

That's it. No flags, no subcommands — just a TUI that walks you through it.

1. **Pick a type** — vCard, Wi-Fi, or URL.
2. **Fill in the fields** — the form shows every field for that type, with a
   live QR preview updating in real time as you type.
3. **Generate** — commit and export to PNG, or tweak the style.

### Keyboard

| Key | Action |
| --- | --- |
| `tab` / `shift+tab` | move between fields |
| `enter` | confirm a field |
| `ctrl+s` | commit and generate the QR |
| `esc` | back / cancel |

On-screen instructions are always visible, so you never have to guess what a
key does.

### vCard first

The vCard flow is the star of the show. It collects the full set of contact
fields — name, phone, email, company, title, address, website, and notes — so
someone can scan your screen and drop you straight into their address book.


## Roadmap

- [ ] QR style menu — colors, gradients, rounded corners, logos
- [ ] vCard template presets
- [ ] Batch export
