package main

import (
	"strings"
	"testing"

	"github.com/skip2/go-qrcode"
)

func setField(fields []field, key, value string) {
	for i := range fields {
		if fields[i].key == key {
			fields[i].input.SetValue(value)
		}
	}
}

func TestVcardContent(t *testing.T) {
	fields := newFields(vcardSpecs)
	setField(fields, "FN", "Jane Doe")
	setField(fields, "TEL", "+1 555 123 4567")
	setField(fields, "EMAIL", "jane@example.com")

	got := vcardContent(fields)
	for _, want := range []string{
		"BEGIN:VCARD",
		"VERSION:3.0",
		"FN:Jane Doe",
		"TEL;TYPE=CELL:+1 555 123 4567",
		"EMAIL:jane@example.com",
		"END:VCARD",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("vcard missing %q:\n%s", want, got)
		}
	}
}

func TestWifiContent(t *testing.T) {
	fields := newFields(wifiSpecs)
	setField(fields, "S", "MyWiFi")
	setField(fields, "P", "hunter2")

	got := wifiContent(fields)
	if got != "WIFI:T:WPA;S:MyWiFi;P:hunter2;;" {
		t.Errorf("wifi content = %q", got)
	}
}

func TestURLContent(t *testing.T) {
	fields := newFields(urlSpecs)
	setField(fields, "URL", "https://example.com")

	if got := urlContent(fields); got != "https://example.com" {
		t.Errorf("url content = %q", got)
	}
}

func TestRenderQR(t *testing.T) {
	qr, err := qrcode.New("hello", qrcode.Medium)
	if err != nil {
		t.Fatal(err)
	}
	out := renderQR(qr)
	if out == "" {
		t.Fatal("renderQR returned empty output")
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("no lines rendered")
	}
	// Half-block rendering packs two module rows per line.
	bm := qr.Bitmap()
	want := (len(bm) + 1) / 2
	if len(lines) != want {
		t.Errorf("expected %d lines for a %dx%d QR, got %d", want, len(bm), len(bm), len(lines))
	}
}

func TestRenderQRStringEmpty(t *testing.T) {
	if got := renderQRString(""); !strings.Contains(got, "fill in fields") {
		t.Errorf("empty content should show a hint, got %q", got)
	}
}

func TestNewFieldsFocus(t *testing.T) {
	fields := newFields(vcardSpecs)
	if len(fields) != len(vcardSpecs) {
		t.Fatalf("expected %d fields, got %d", len(vcardSpecs), len(fields))
	}
	if !fields[0].input.Focused() {
		t.Error("first field should be focused")
	}
}
