package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	c, err := parseHexColor("#FF6599")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.R != 0xFF || c.G != 0x65 || c.B != 0x99 || c.A != 255 {
		t.Errorf("expected RGBA(255, 101, 153, 255), got %+v", c)
	}

	c2, err := parseHexColor("784BA0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c2.R != 0x78 || c2.G != 0x4B || c2.B != 0xA0 {
		t.Errorf("expected RGBA(120, 75, 160, 255), got %+v", c2)
	}

	if _, err := parseHexColor("#12"); err == nil {
		t.Error("expected error for short hex")
	}
	if _, err := parseHexColor("ZZZZZZ"); err == nil {
		t.Error("expected error for invalid hex chars")
	}
}

func TestInterpolateColor(t *testing.T) {
	c1 := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	c2 := color.RGBA{R: 200, G: 100, B: 50, A: 255}

	start := interpolateColor(c1, c2, 0.0)
	if start != c1 {
		t.Errorf("t=0: got %+v, want %+v", start, c1)
	}

	end := interpolateColor(c1, c2, 1.0)
	if end != c2 {
		t.Errorf("t=1: got %+v, want %+v", end, c2)
	}

	mid := interpolateColor(c1, c2, 0.5)
	if mid.R != 100 || mid.G != 50 || mid.B != 25 {
		t.Errorf("t=0.5: got %+v, want (100, 50, 25)", mid)
	}
}

func TestColorToHex(t *testing.T) {
	c := color.RGBA{R: 255, G: 101, B: 153, A: 255}
	if hex := colorToHex(c); hex != "#FF6599" {
		t.Errorf("got %q, want #FF6599", hex)
	}
}

func TestPresetsValidity(t *testing.T) {
	if len(presets) < 5 {
		t.Fatalf("expected at least 5 presets, got %d", len(presets))
	}
	for _, p := range presets {
		if _, err := parseHexColor(p.FgStart); err != nil {
			t.Errorf("preset %s has invalid FgStart %q: %v", p.Name, p.FgStart, err)
		}
		if _, err := parseHexColor(p.FgEnd); err != nil {
			t.Errorf("preset %s has invalid FgEnd %q: %v", p.Name, p.FgEnd, err)
		}
		if _, err := parseHexColor(p.BgColor); err != nil {
			t.Errorf("preset %s has invalid BgColor %q: %v", p.Name, p.BgColor, err)
		}
	}
}

func TestRenderStyledPNG(t *testing.T) {
	testFile := "test-styled.png"
	defer os.Remove(testFile)

	style := QRStyle{
		PresetIndex: 1, // Giant Goldfish
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
	}

	err := renderStyledPNG("https://example.com/styled", style, 256, testFile)
	if err != nil {
		t.Fatalf("renderStyledPNG failed: %v", err)
	}
	if _, err := os.Stat(testFile); err != nil {
		t.Fatalf("expected %s to exist: %v", testFile, err)
	}
}

func TestRenderStyledPNGWithLogo(t *testing.T) {
	logoFile := "test-logo.png"
	outFile := "test-styled-logo.png"
	defer os.Remove(logoFile)
	defer os.Remove(outFile)

	// Create a dummy 32x32 logo PNG
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 255, G: 0, B: 128, A: 255})
		}
	}
	f, err := os.Create(logoFile)
	if err != nil {
		t.Fatalf("failed to create logo file: %v", err)
	}
	png.Encode(f, img)
	f.Close()

	style := QRStyle{
		PresetIndex: 0, // Candy Glam
		Shape:       ShapeRounded,
		Gradient:    GradientVertical,
		LogoPath:    logoFile,
	}

	err = renderStyledPNG("https://example.com/logo-test", style, 256, outFile)
	if err != nil {
		t.Fatalf("renderStyledPNG with logo failed: %v", err)
	}
	if _, err := os.Stat(outFile); err != nil {
		t.Fatalf("expected %s to exist: %v", outFile, err)
	}
}

func TestRenderStyledQR(t *testing.T) {
	style := QRStyle{
		PresetIndex: 0,
		Shape:       ShapeRounded,
		Gradient:    GradientHorizontal,
	}

	out := renderStyledQRString("https://example.com", style)
	if out == "" {
		t.Fatal("expected non-empty styled QR string")
	}
	if !strings.Contains(out, "▀") && !strings.Contains(out, "▄") {
		t.Errorf("expected half-block characters in styled QR output: %s", out)
	}
}
