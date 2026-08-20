package main

import (
	"fmt"
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

func vcardContent(fields []field) string {
	name := ""
	for _, f := range fields {
		if f.key == "FN" {
			name = strings.TrimSpace(f.input.Value())
		}
	}
	if name == "" {
		return ""
	}

	var b strings.Builder
	b.WriteString("BEGIN:VCARD\nVERSION:3.0\n")
	for _, f := range fields {
		v := strings.TrimSpace(f.input.Value())
		if v == "" {
			continue
		}
		switch f.key {
		case "FN":
			b.WriteString("FN:" + v + "\n")
		case "TEL":
			b.WriteString("TEL;TYPE=CELL:" + v + "\n")
		case "EMAIL":
			b.WriteString("EMAIL:" + v + "\n")
		case "ORG":
			b.WriteString("ORG:" + v + "\n")
		case "TITLE":
			b.WriteString("TITLE:" + v + "\n")
		case "ADR":
			b.WriteString("ADR;TYPE=WORK:;;" + v + ";;;;\n")
		case "URL":
			b.WriteString("URL:" + v + "\n")
		case "NOTE":
			b.WriteString("NOTE:" + v + "\n")
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
	return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", enc, ssid, password)
}

func urlContent(fields []field) string {
	if len(fields) == 0 {
		return ""
	}
	return fields[0].input.Value()
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
	if err := qrcode.WriteFile(content, qrcode.Medium, 512, "candy-qr.png"); err != nil {
		m.message = "Export failed: " + err.Error()
		return
	}
	m.message = "Saved to candy-qr.png"
}
