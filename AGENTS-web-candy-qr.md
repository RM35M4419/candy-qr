# AGENTS-web-candy-qr.md

Guide for AI agents (and humans) building and maintaining the Web edition of **Candy QR**.

---

## 1. What this is

**Candy QR Web** is a single-page web application companion to the Linux TUI app. It generates, styles, and exports beautiful, scannable QR codes in real time. It covers three encodings (**vCard 3.0**, **Wi-Fi**, and **URL**), provides curated [COLOURlovers](https://www.colourlovers.com/) gradient palettes, rounded module shapes, and center logo embedding.

---

## 2. Stack & Constraints

- **Core Technologies:** HTML5, Vanilla CSS3, Vanilla ES6+ JavaScript.
- **No Heavy Frameworks:** No React, Vue, Angular, or Tailwind build steps. Keep it light, instant to load, and deployable as static assets on GitHub Pages, Cloudflare Pages, or Netlify.
- **QR Encoding Engine:** Lightweight zero-dependency QR matrix generator (e.g. embedded `qrcode-generator` or custom JS matrix calculation) outputting boolean module grids.
- **Rendering Engine:** HTML5 `<canvas>` API with custom gradient shaders, subpixel rounded rectangle drawers, and high-fidelity image compositing.
- **Export:** Canvas `toBlob` / `toDataURL` for high-DPI 1024×1024 PNG and dynamic SVG markup generator.

---

## 3. Aesthetics & Brand Voice (Charm.land Style)

Follow the design language of [charm.land](https://charm.land/):
- **Surfaces:** Dark obsidian background (`#0B0C10`), elevated panels (`#14161D`), subtle border stroke (`#27272A`), soft neon outer glow on active inputs.
- **Typography:**
  - Code/Headings/Badges: `'JetBrains Mono', 'Fira Code', ui-monospace, monospace`
  - Body/Labels: `'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`
- **Keycap Badges (`<kbd>`):** Retro terminal style with border and shadow.
- **Interactive Feedback:** Smooth CSS transitions (`cubic-bezier(0.16, 1, 0.3, 1)`), hover micro-animations, active button presses, and subtle pulse indicators.
- **Copy:** Playful, warm, benefit-driven (Charm brand voice). Keep whimsy in UI copy, keep logic clean and robust.

---

## 4. Architecture & File Structure

```
candy-qr-web/
├── index.html         # Semantic structure, layout grid, modals/drawers
├── css/
│   ├── reset.css      # Modern box-sizing, baseline resets
│   ├── tokens.css     # CSS Variables (palettes, spacing, fonts, radiuses)
│   └── style.css      # Layout, Charm aesthetic components, responsiveness
├── js/
│   ├── app.js         # State coordinator & DOM event listeners
│   ├── qr-matrix.js   # QR algorithm matrix generator (1..40, L/M/Q/H)
│   ├── encoder.js     # vCard, Wi-Fi, URL payload formatters & escaping
│   ├── renderer.js    # Canvas rendering (gradients, rounded shapes, logo)
│   ├── presets.js     # 10 COLOURlovers preset definitions (5 light + 5 dark)
│   └── export.js      # PNG download, SVG download, Clipboard copy
└── assets/
    ├── icons/         # SVG icons (vCard, Wi-Fi, URL, Download, Copy, Logo)
    └── favicon.png    # App favicon
```

---

## 5. State Management & Flow

Use a simple unidirectional reactive state model:

```js
const state = {
  activeType: 'vcard', // 'vcard' | 'wifi' | 'url'
  fields: {
    vcard: { given: '', family: '', tel: '', tel2: '', email: '', email2: '', org: '', title: '', adr: '', url: '', note: '' },
    wifi: { ssid: '', password: '', encryption: 'WPA' },
    url: { url: 'https://' }
  },
  style: {
    presetIndex: 0,
    shape: 'rounded',     // 'rounded' | 'square'
    gradient: 'diagonal', // 'diagonal' | 'vertical' | 'horizontal'
    logoImage: null,      // Image object or null
    logoFileName: ''
  }
};
```

Whenever any input changes:
1. `encoder.buildPayload(state)` computes the encoded string.
2. `qr-matrix.generate(payload, errorCorrectionLevel)` computes the `boolean[][]` grid.
3. `renderer.draw(canvas, matrix, state.style)` draws the live preview canvas.

---

## 6. QR Payload Encodings (Strict Compliance)

### 1. vCard 3.0 (RFC 2426)
- Escape special characters: `\` ➔ `\\`, `;` ➔ `\;`, `,` ➔ `\,`, newlines ➔ `\n`.
- If First and Last name are omitted, fallback `FN` to Company (`ORG`), `EMAIL`, or `TEL` (default to `"Contact"`).
- Structure:
  ```
  BEGIN:VCARD
  VERSION:3.0
  N:Family;Given;;;
  FN:Given Family
  TEL;TYPE=CELL:+1 555 123 4567
  EMAIL:name@example.com
  ORG:Company Name
  TITLE:Job Title
  ADR;TYPE=WORK:;;Street Address;;;;
  URL:https://example.com
  NOTE:Notes
  END:VCARD
  ```

### 2. Wi-Fi (ZXing Standard)
- Format: `WIFI:T:<enc>;S:<ssid>;P:<password>;;`
- Escape special characters in SSID & Password: `\`, `;`, `:`, `,`, `"`.
- If SSID is blank, return empty string (do not render invalid empty QR).

### 3. URL
- Trim whitespace, ensure valid URI format.

---

## 7. Curated Color Presets (COLOURlovers)

Provide exactly **10 presets (5 Light + 5 Dark)**:

```js
const PRESETS = [
  // 5 Light Backgrounds
  { id: 'candy',    name: 'Candy Glam 🍬',       fg1: '#FF6599', fg2: '#784BA0', bg: '#FFFFFF', shape: 'rounded', gradient: 'diagonal' },
  { id: 'goldfish', name: 'Giant Goldfish 🐠',    fg1: '#FA6900', fg2: '#69D2E7', bg: '#FFFFFF', shape: 'rounded', gradient: 'diagonal' },
  { id: 'emo',      name: 'Cheer Up 🎈',         fg1: '#FF6B6B', fg2: '#4ECDC4', bg: '#FFFFFF', shape: 'rounded', gradient: 'diagonal' },
  { id: 'ocean',    name: 'Ocean Five 🌊',       fg1: '#00A0B0', fg2: '#EB6841', bg: '#FAF7F2', shape: 'rounded', gradient: 'diagonal' },
  { id: 'thought',  name: 'Thought Provoking 🍷', fg1: '#D95B43', fg2: '#542437', bg: '#FDF8F5', shape: 'rounded', gradient: 'diagonal' },

  // 5 Dark / Black Backgrounds
  { id: 'cyber',    name: 'Cyber Neon ⚡',        fg1: '#00F5D4', fg2: '#FF007F', bg: '#0B0C10', shape: 'rounded', gradient: 'diagonal' },
  { id: 'synthwave',name: 'Synthwave 80s 🕶️',    fg1: '#7209B7', fg2: '#4CC9F0', bg: '#05050A', shape: 'rounded', gradient: 'diagonal' },
  { id: 'magma',    name: 'Firewatch Magma 🔥',   fg1: '#FF4E50', fg2: '#F9D423', bg: '#0D0D11', shape: 'rounded', gradient: 'diagonal' },
  { id: 'toxic',    name: 'Toxic Lime 🧪',        fg1: '#C7F464', fg2: '#00F5D4', bg: '#0D1117', shape: 'rounded', gradient: 'diagonal' },
  { id: 'darkmono', name: 'Dark Monochrome 🌘',   fg1: '#FFFFFF', fg2: '#A0AEC0', bg: '#000000', shape: 'square',  gradient: 'diagonal' },
];
```

---

## 8. Canvas Rendering & Logo Insertion Rules

1. **Error Correction Level:**
   - Standard: `Level Q` or `Level H`.
   - When center logo is present: **Always enforce `Level H` (30% error recovery)**.
2. **Gradient Factor Calculation:**
   - Diagonal: `t = (x + y) / (2 * N - 2)`
   - Vertical: `t = y / (N - 1)`
   - Horizontal: `t = x / (N - 1)`
3. **Rounded Module Corners:**
   - Draw rounded rectangles using `ctx.roundRect(x, y, w, h, radius)` or arc paths with `radius = modSize * 0.35`.
4. **Center Logo Plate:**
   - Badge size: ~22–24% of total canvas size.
   - Draw background badge in `bg` color with rounded corners (`badgeRadius = badgeSize * 0.20`) to separate modules from the logo.
   - Scale logo with aspect-ratio preservation and anti-aliasing (`ctx.imageSmoothingQuality = 'high'`).

---

## 9. Scope Guardrails

- **Keep It Simple & Focused:** This is a dedicated, fast QR generator. Do not bloat with backend databases, authentication, or advertising scripts.
- **Client-Side Privacy First:** All QR generation and logo processing happens 100% in the user's browser canvas. No contact data or logos ever leave the client.
- **Accessibility:** Full keyboard navigation (`tabindex`, visible focus rings), screen reader labels (`aria-label`), and color contrast.
