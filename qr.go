package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/skip2/go-qrcode"
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
	return "Preview\n\n" + renderQRString(m.buildContent())
}

func renderQRString(content string) string {
	if strings.TrimSpace(content) == "" {
		return "(fill in fields to see the QR)"
	}
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "(error generating QR)"
	}
	return renderQR(qr)
}

func renderQR(q *qrcode.QRCode) string {
	bm := q.Bitmap()
	var b strings.Builder
	for y := 0; y < len(bm); y += 2 {
		for x := 0; x < len(bm[y]); x++ {
			top := bm[y][x]
			bottom := false
			if y+1 < len(bm) {
				bottom = bm[y+1][x]
			}
			switch {
			case top && bottom:
				b.WriteString("█")
			case top && !bottom:
				b.WriteString("▀")
			case !top && bottom:
				b.WriteString("▄")
			default:
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *Model) exportPNG() {
	content := m.buildContent()
	if strings.TrimSpace(content) == "" {
		m.message = "Nothing to export yet"
		return
	}
	filename := "candy-qr.png"
	if err := qrcode.WriteFile(content, qrcode.Medium, 512, filename); err != nil {
		m.message = "Export failed: " + err.Error()
		return
	}
	if abs, err := filepath.Abs(filename); err == nil {
		m.message = "Saved to " + abs
	} else {
		m.message = "Saved to " + filename
	}
}
