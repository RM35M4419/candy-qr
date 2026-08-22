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

func TestVcardContentMultipleContacts(t *testing.T) {
	fields := newFields(vcardSpecs)
	setField(fields, "GIVEN", "Jane")
	setField(fields, "FAMILY", "Doe")
	setField(fields, "TEL", "111")
	setField(fields, "TEL2", "222")
	setField(fields, "EMAIL", "a@b.com")
	setField(fields, "EMAIL2", "c@d.com")

	got := vcardContent(fields)
	if n := strings.Count(got, "TEL;TYPE=CELL:"); n != 2 {
		t.Errorf("expected 2 phone numbers, got %d:\n%s", n, got)
	}
	if n := strings.Count(got, "EMAIL:"); n != 2 {
		t.Errorf("expected 2 emails, got %d:\n%s", n, got)
	}
	for _, want := range []string{"TEL;TYPE=CELL:111", "TEL;TYPE=CELL:222", "EMAIL:a@b.com", "EMAIL:c@d.com"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestVcardContent(t *testing.T) {
	fields := newFields(vcardSpecs)
	setField(fields, "GIVEN", "Jane")
	setField(fields, "FAMILY", "Doe")
	setField(fields, "TEL", "+1 555 123 4567")
	setField(fields, "EMAIL", "jane@example.com")

	got := vcardContent(fields)
	for _, want := range []string{
		"BEGIN:VCARD",
		"VERSION:3.0",
		"N:Doe;Jane;;;",
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

func TestVcardContentEscaping(t *testing.T) {
	fields := newFields(vcardSpecs)
	setField(fields, "GIVEN", "Jane;Jr.")
	setField(fields, "FAMILY", "Doe, Esq.")
	setField(fields, "ORG", "Acme, Inc.; Division\\1")
	setField(fields, "NOTE", "Special, note; with \\ backslash")

	got := vcardContent(fields)
	for _, want := range []string{
		"N:Doe\\, Esq.;Jane\\;Jr.;;;",
		"FN:Jane\\;Jr. Doe\\, Esq.",
		"ORG:Acme\\, Inc.\\; Division\\\\1",
		"NOTE:Special\\, note\\; with \\\\ backslash",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("vcard missing escaped %q in:\n%s", want, got)
		}
	}

	if gotEsc := escapeVcard("Line 1\nLine 2; text, 1\\2"); gotEsc != "Line 1\\nLine 2\\; text\\, 1\\\\2" {
		t.Errorf("escapeVcard mismatch: got %q", gotEsc)
	}
}

func TestVcardContentNamelessFallback(t *testing.T) {
	fields := newFields(vcardSpecs)
	setField(fields, "ORG", "Acme Corp")
	setField(fields, "EMAIL", "contact@acme.com")

	got := vcardContent(fields)
	if got == "" {
		t.Fatal("expected non-empty vcard for company contact")
	}
	if !strings.Contains(got, "FN:Acme Corp") {
		t.Errorf("expected FN:Acme Corp, got:\n%s", got)
	}
}

func TestVcardContentEmpty(t *testing.T) {
	fields := newFields(vcardSpecs)
	if got := vcardContent(fields); got != "" {
		t.Errorf("expected empty string for blank vcard fields, got %q", got)
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

func TestWifiContentEmpty(t *testing.T) {
	fields := newFields(wifiSpecs)
	if got := wifiContent(fields); got != "" {
		t.Errorf("expected empty string for blank wifi SSID, got %q", got)
	}
}

func TestWifiContentEscaping(t *testing.T) {
	fields := newFields(wifiSpecs)
	setField(fields, "S", `My;Wi-Fi:5G\Home,"Office"`)
	setField(fields, "P", `pass;word:123,"xyz"\`)

	got := wifiContent(fields)
	want := `WIFI:T:WPA;S:My\;Wi-Fi\:5G\\Home\,\"Office\";P:pass\;word\:123\,\"xyz\"\\;;`
	if got != want {
		t.Errorf("wifi content =\n%q\nwant:\n%q", got, want)
	}
}

func TestURLContent(t *testing.T) {
	fields := newFields(urlSpecs)
	setField(fields, "URL", "  https://example.com  ")

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
