# 🍬 Candy QR Web

> Your browser, but make it a delicious QR code business card.

**Candy QR Web** is a fast, zero-dependency browser app for making beautiful, scannable QR codes. Build a vCard contact, Wi-Fi credentials, or a clean link, style it with curated [COLOURlovers](https://www.colourlovers.com/) palettes and gradients, drop your logo in the center, and export a crisp, high-res PNG or SVG.

Designed with the vibrant, terminal-inspired charm of [charm.land](https://charm.land/) — sleek dark surfaces, glowing neon accents, keyboard ergonomics, and delightful micro-interactions.

---

## ✨ Features

- **📇 vCard First**: Full contact card encoding (Name, Phone, Email, Company, Title, Address, Website, Notes) that drops straight into iOS and Android address books upon scan.
- **📶 Wi-Fi Sharing**: Connect guests instantly with encrypted `WPA`/`WEP`/`nopass` credentials.
- **🔗 Smart URL**: Turn any link into an instant scannable destination.
- **🎨 10 Curated Palettes (5 Light + 5 Dark)**: Hand-picked top palettes from COLOURlovers (*Candy Glam*, *Giant Goldfish*, *Cheer Up*, *Cyber Neon*, *Synthwave 80s*, *Firewatch Magma*, and more).
- **🌈 Gradients & Shapes**: Diagonal, vertical, or horizontal gradients with smooth rounded pill modules or classic square blocks.
- **🖼️ Center Logo Embedding**: Drag and drop any PNG/SVG logo into the center. Automatic high-redundancy error correction (`Level H` — 30% recovery) keeps codes 100% scannable.
- **⚡ Live Reactive Canvas**: Real-time rendering as you type with subpixel anti-aliasing.
- **📥 1-Click Export**: Save high-resolution 1024×1024 PNGs, download SVGs, or copy directly to clipboard.
- **⌨️ Keyboard Ergonomics**: Full tab navigation, enter confirms, and shortcut keys.

---

## 🛠️ Stack

- **Structure:** Semantic HTML5
- **Styling:** Vanilla CSS (CSS Variables, Flexbox/Grid, Glassmorphism, Charm-inspired typography)
- **Logic:** Vanilla ES6+ JavaScript (zero external runtime dependencies, lightweight client-side QR generation via HTML5 `<canvas>`)
- **Compatibility:** All modern browsers (Chrome, Firefox, Safari, Edge) on desktop and mobile.

---

## 🚀 Getting Started

No build tools, bundlers, or `node_modules` required!

```bash
# Clone the repository
git clone https://github.com/RM35M4419/candy-qr-web.git
cd candy-qr-web

# Serve locally with any static web server
python3 -m http.server 8000
# or: npx serve .
```

Open `http://localhost:8000` in your browser.

---

## 🎯 Usage Flow

1. **Pick a Type:** Select **vCard**, **Wi-Fi**, or **URL**.
2. **Fill the Fields:** Enter your details. Watch the live QR update beside your form in real time.
3. **Customize Style:**
   - Pick from **10 curated palettes** (5 light themes, 5 dark themes).
   - Toggle **Rounded** or **Square** modules.
   - Choose **Diagonal**, **Vertical**, or **Horizontal** gradients.
   - Drop a **Center Logo** (PNG / SVG).
4. **Export & Share:**
   - Click **Download PNG (1024×1024)**.
   - Click **Download SVG**.
   - Or click **Copy to Clipboard**.

---

## ⌨️ Keyboard Shortcuts

| Shortcut | Action |
| --- | --- |
| `Tab` / `Shift+Tab` | Navigate through form and style controls |
| `Enter` | Advance to next field or trigger button |
| `Ctrl+S` / `Cmd+S` | Download high-resolution PNG |
| `Ctrl+C` / `Cmd+C` | Copy QR image to clipboard (when preview focused) |
| `Esc` | Reset or clear active modal/popover |

---

## 🎨 Charm Design System

Candy QR Web is crafted following the aesthetic principles of **Charm.land**:
- **Backgrounds:** Deep obsidian `#0B0C10` and slate `#16181D`.
- **Card Surfaces:** Semi-transparent dark glass with subtle border outlines (`#27272A`).
- **Accent Glows:** Neon Cyan (`#00F5D4`), Hot Magenta (`#FF007F`), and Electric Purple (`#7B2CBF`).
- **Typography:** Clean blend of monospace headers/badges (`JetBrains Mono`, `Fira Code`) with modern sans-serif body (`Inter`, `system-ui`).

---

## 📄 License

MIT © Candy QR Team
