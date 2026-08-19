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

- **Wi-Fi** — share your network name and password as a scannable code
- **URL** — turn any link into a QR code
- **vCard** — encode your name, phone, email, and more directly into a
  scannable contact that drops into a phone's address book
- **Print to terminal** — render the QR in your terminal (when it supports it)
- **Export to PNG** — save a crisp, shareable image

## Status

**v0.1** — the delicious-but-plain first cut. Functionality first, no colors
or fancy design _yet_. Logo, gradients, rounded corners, and the full glam
treatment are on the roadmap.

## Install

```bash
go install github.com/yourname/candy-qr@latest
```

## Usage

```
candy-qr                 # start the interactive TUI
candy-qr url <url>       # make a QR from a URL
candy-qr wifi <ssid> <password>
candy-qr vcard <file>    # make a QR from a vCard file
```

Run it, pick what you want to share, fill in the fields, and done. Print it
to the terminal or export a PNG with a keystroke.

## Why?

Your contact info, Wi-Fi password, and links all have one thing in common:
people type them out by hand, and miss a character. QR codes mean you hold up
a screen and the other person's camera does the rest. Candy QR is the pretty
little tool that lives in that gap.

## Roadmap

- [ ] Logo, colors, and gradients
- [ ] Rounded corners and styling options
- [ ] vCard template presets
- [ ] Batch export

## Contributing

Pull requests welcome. Keep it sweet.

## License

[MIT](LICENSE)