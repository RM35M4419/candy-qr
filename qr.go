package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/skip2/go-qrcode"
	xdraw "golang.org/x/image/draw"
)

func (m Model) buildContent() string {
	switch qrTypes[m.typeCursor] {
	case "vCard":
		return vcardContent(m.fields)
	case "Wi-Fi":
		return wifiContent(m.fields)
	case "URL":
		return urlContent(m.fields)
	}
	return ""
}

func escapeVcard(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\r\n", "\\n")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

func escapeWifi(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\', ';', ',', ':', '"':
			b.WriteRune('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func vcardContent(fields []field) string {
	given, family := "", ""
	org, email, tel := "", "", ""
	hasAnyField := false

	for _, f := range fields {
		v := strings.TrimSpace(f.input.Value())
		if v == "" {
			continue
		}
		hasAnyField = true
		switch f.key {
		case "GIVEN":
			given = v
		case "FAMILY":
			family = v
		case "ORG":
			org = v
		case "EMAIL", "EMAIL2":
			if email == "" {
				email = v
			}
		case "TEL", "TEL2":
			if tel == "" {
				tel = v
			}
		}
	}

	if !hasAnyField {
		return ""
	}

	formatted := strings.TrimSpace(given + " " + family)
	if formatted == "" {
		if org != "" {
			formatted = org
		} else if email != "" {
			formatted = email
		} else if tel != "" {
			formatted = tel
		} else {
			formatted = "Contact"
		}
	}

	var b strings.Builder
	b.WriteString("BEGIN:VCARD\nVERSION:3.0\n")
	if given != "" || family != "" {
		b.WriteString("N:" + escapeVcard(family) + ";" + escapeVcard(given) + ";;;\n")
	}
	b.WriteString("FN:" + escapeVcard(formatted) + "\n")
	for _, f := range fields {
		v := strings.TrimSpace(f.input.Value())
		if v == "" {
			continue
		}
		switch f.key {
		case "TEL", "TEL2":
			b.WriteString("TEL;TYPE=CELL:" + escapeVcard(v) + "\n")
		case "EMAIL", "EMAIL2":
			b.WriteString("EMAIL:" + escapeVcard(v) + "\n")
		case "ORG":
			b.WriteString("ORG:" + escapeVcard(v) + "\n")
		case "TITLE":
			b.WriteString("TITLE:" + escapeVcard(v) + "\n")
		case "ADR":
			b.WriteString("ADR;TYPE=WORK:;;" + escapeVcard(v) + ";;;;\n")
		case "URL":
			b.WriteString("URL:" + v + "\n")
		case "NOTE":
			b.WriteString("NOTE:" + escapeVcard(v) + "\n")
		}
	}
	b.WriteString("END:VCARD")
	return b.String()
}

func wifiContent(fields []field) string {
	ssid := ""
	password := ""
	enc := "WPA"
	for _, f := range fields {
		switch f.key {
		case "S":
			ssid = f.input.Value()
		case "P":
			password = f.input.Value()
		case "T":
			if v := strings.TrimSpace(f.input.Value()); v != "" {
				enc = v
			}
		}
	}
	if strings.TrimSpace(ssid) == "" {
		return ""
	}
	return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", escapeWifi(enc), escapeWifi(ssid), escapeWifi(password))
}

func urlContent(fields []field) string {
	if len(fields) == 0 {
		return ""
	}
	return strings.TrimSpace(fields[0].input.Value())
}

func (m Model) renderPreview() string {
	return "Preview\n\n" + renderStyledQRString(m.buildContent(), m.style)
}

func renderQRString(content string) string {
	return renderStyledQRString(content, defaultStyle())
}

func renderStyledQRString(content string, style QRStyle) string {
	if strings.TrimSpace(content) == "" {
		return "(fill in fields to see the QR)"
	}
	level := qrcode.Medium
	if style.LogoPath != "" {
		level = qrcode.Highest
	}
	qr, err := qrcode.New(content, level)
	if err != nil {
		return "(error generating QR)"
	}
	return renderStyledQR(qr, style)
}

func renderQR(q *qrcode.QRCode) string {
	return renderStyledQR(q, defaultStyle())
}

func renderStyledQR(q *qrcode.QRCode, style QRStyle) string {
	bm := q.Bitmap()
	rows := len(bm)
	if rows == 0 {
		return ""
	}
	cols := len(bm[0])

	preset := style.CurrentPreset()
	fg1, err1 := parseHexColor(preset.FgStart)
	fg2, err2 := parseHexColor(preset.FgEnd)
	if err1 != nil {
		fg1 = color.RGBA{R: 255, G: 101, B: 153, A: 255}
	}
	if err2 != nil {
		fg2 = color.RGBA{R: 120, G: 75, B: 160, A: 255}
	}

	getFactor := func(x, y int) float64 {
		switch style.Gradient {
		case GradientVertical:
			if rows <= 1 {
				return 0
			}
			return float64(y) / float64(rows-1)
		case GradientHorizontal:
			if cols <= 1 {
				return 0
			}
			return float64(x) / float64(cols-1)
		default:
			if rows+cols <= 2 {
				return 0
			}
			return float64(x+y) / float64(rows+cols-2)
		}
	}

	var b strings.Builder
	for y := 0; y < rows; y += 2 {
		for x := 0; x < cols; x++ {
			top := bm[y][x]
			bottom := false
			if y+1 < rows {
				bottom = bm[y+1][x]
			}

			tTop := getFactor(x, y)
			tBot := getFactor(x, y+1)
			topHex := colorToHex(interpolateColor(fg1, fg2, tTop))
			botHex := colorToHex(interpolateColor(fg1, fg2, tBot))

			switch {
			case top && bottom:
				b.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color(topHex)).
					Background(lipgloss.Color(botHex)).
					Render("▀"))
			case top && !bottom:
				b.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color(topHex)).
					Background(lipgloss.Color(preset.BgColor)).
					Render("▀"))
			case !top && bottom:
				b.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color(botHex)).
					Background(lipgloss.Color(preset.BgColor)).
					Render("▄"))
			default:
				b.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color(preset.BgColor)).
					Background(lipgloss.Color(preset.BgColor)).
					Render(" "))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func drawRoundedRect(img *image.RGBA, x0, y0, x1, y1 int, radius float64, c color.RGBA) {
	if radius <= 0 {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				img.SetRGBA(x, y, c)
			}
		}
		return
	}
	cx0 := float64(x0) + radius
	cy0 := float64(y0) + radius
	cx1 := float64(x1) - radius
	cy1 := float64(y1) - radius

	for y := y0; y < y1; y++ {
		fy := float64(y) + 0.5
		for x := x0; x < x1; x++ {
			fx := float64(x) + 0.5
			var dist float64 = -1.0
			if fx < cx0 && fy < cy0 {
				dx, dy := fx-cx0, fy-cy0
				dist = math.Hypot(dx, dy) - radius
			} else if fx > cx1 && fy < cy0 {
				dx, dy := fx-cx1, fy-cy0
				dist = math.Hypot(dx, dy) - radius
			} else if fx < cx0 && fy > cy1 {
				dx, dy := fx-cx0, fy-cy1
				dist = math.Hypot(dx, dy) - radius
			} else if fx > cx1 && fy > cy1 {
				dx, dy := fx-cx1, fy-cy1
				dist = math.Hypot(dx, dy) - radius
			}

			if dist <= -0.5 {
				img.SetRGBA(x, y, c)
			} else if dist < 0.5 {
				alphaFactor := 0.5 - dist
				if alphaFactor > 0 {
					orig := img.RGBAAt(x, y)
					a := float64(c.A) / 255.0 * alphaFactor
					nr := uint8(float64(c.R)*a + float64(orig.R)*(1.0-a) + 0.5)
					ng := uint8(float64(c.G)*a + float64(orig.G)*(1.0-a) + 0.5)
					nb := uint8(float64(c.B)*a + float64(orig.B)*(1.0-a) + 0.5)
					img.SetRGBA(x, y, color.RGBA{R: nr, G: ng, B: nb, A: 255})
				}
			}
		}
	}
}

func scaleAndOverlayLogo(dst *image.RGBA, logo image.Image, x0, y0, x1, y1 int) {
	targetW := x1 - x0
	targetH := y1 - y0
	if targetW <= 0 || targetH <= 0 {
		return
	}
	bounds := logo.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW == 0 || srcH == 0 {
		return
	}

	scale := math.Min(float64(targetW)/float64(srcW), float64(targetH)/float64(srcH))
	w := int(float64(srcW)*scale + 0.5)
	h := int(float64(srcH)*scale + 0.5)
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	offsetX := x0 + (targetW-w)/2
	offsetY := y0 + (targetH-h)/2

	dstRect := image.Rect(offsetX, offsetY, offsetX+w, offsetY+h)
	// Catmull-Rom high-quality bicubic downsampling filter
	xdraw.CatmullRom.Scale(dst, dstRect, logo, bounds, xdraw.Over, nil)
}

func renderStyledPNG(content string, style QRStyle, size int, filename string) error {
	level := qrcode.High
	if style.LogoPath != "" {
		level = qrcode.Highest
	}

	qr, err := qrcode.New(content, level)
	if err != nil {
		return err
	}

	bm := qr.Bitmap()
	n := len(bm)
	if n == 0 {
		return fmt.Errorf("empty QR matrix")
	}

	modSize := size / n
	if modSize < 1 {
		modSize = 1
	}
	imgSize := n * modSize

	preset := style.CurrentPreset()
	fg1, err1 := parseHexColor(preset.FgStart)
	fg2, err2 := parseHexColor(preset.FgEnd)
	bg, err3 := parseHexColor(preset.BgColor)
	if err1 != nil {
		fg1 = color.RGBA{R: 255, G: 101, B: 153, A: 255}
	}
	if err2 != nil {
		fg2 = color.RGBA{R: 120, G: 75, B: 160, A: 255}
	}
	if err3 != nil {
		bg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}

	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	for y := 0; y < imgSize; y++ {
		for x := 0; x < imgSize; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	getFactor := func(x, y int) float64 {
		switch style.Gradient {
		case GradientVertical:
			if n <= 1 {
				return 0
			}
			return float64(y) / float64(n-1)
		case GradientHorizontal:
			if n <= 1 {
				return 0
			}
			return float64(x) / float64(n-1)
		default:
			if 2*n <= 2 {
				return 0
			}
			return float64(x+y) / float64(2*n-2)
		}
	}

	radius := 0.0
	if style.Shape == ShapeRounded {
		radius = float64(modSize) * 0.35
	}

	for my := 0; my < n; my++ {
		for mx := 0; mx < n; mx++ {
			if bm[my][mx] {
				t := getFactor(mx, my)
				modColor := interpolateColor(fg1, fg2, t)
				x0 := mx * modSize
				y0 := my * modSize
				x1 := x0 + modSize
				y1 := y0 + modSize
				drawRoundedRect(img, x0, y0, x1, y1, radius, modColor)
			}
		}
	}

	if style.LogoPath != "" {
		logoFile, err := os.Open(style.LogoPath)
		if err == nil {
			defer logoFile.Close()
			logoImg, _, err := image.Decode(logoFile)
			if err == nil {
				badgeSize := int(float64(imgSize) * 0.24)
				cx := imgSize / 2
				cy := imgSize / 2
				bx0 := cx - badgeSize/2
				by0 := cy - badgeSize/2
				bx1 := bx0 + badgeSize
				by1 := by0 + badgeSize

				badgeRadius := float64(badgeSize) * 0.20
				drawRoundedRect(img, bx0, by0, bx1, by1, badgeRadius, bg)

				pad := int(float64(badgeSize) * 0.10)
				scaleAndOverlayLogo(img, logoImg, bx0+pad, by0+pad, bx1-pad, by1-pad)
			}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func (m *Model) exportPNG() {
	content := m.buildContent()
	if strings.TrimSpace(content) == "" {
		m.message = "Nothing to export yet"
		return
	}
	filename := "candy-qr.png"
	if err := renderStyledPNG(content, m.style, 1024, filename); err != nil {
		m.message = "Export failed: " + err.Error()
		return
	}
	if abs, err := filepath.Abs(filename); err == nil {
		m.message = "Saved to " + abs
	} else {
		m.message = "Saved to " + filename
	}
}
