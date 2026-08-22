package main

import (
	"fmt"
	"image/color"
	"strconv"
)

type GradientDirection int

const (
	GradientDiagonal GradientDirection = iota
	GradientVertical
	GradientHorizontal
)

var gradientNames = []string{"Diagonal ↘", "Vertical ↓", "Horizontal →"}

type ModuleShape int

const (
	ShapeRounded ModuleShape = iota
	ShapeSquare
)

var shapeNames = []string{"Rounded", "Square"}

type StylePreset struct {
	ID          string
	Name        string
	Description string
	FgStart     string
	FgEnd       string
	BgColor     string
	Shape       ModuleShape
	Gradient    GradientDirection
}

var presets = []StylePreset{
	{
		ID:          "candy",
		Name:        "Candy Glam 🍬",
		Description: "Vibrant candy pink to royal purple gradient",
		FgStart:     "#FF6599",
		FgEnd:       "#784BA0",
		BgColor:     "#FFFFFF",
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
	},
	{
		ID:          "goldfish",
		Name:        "Giant Goldfish 🐠",
		Description: "COLOURlovers #1 favorite: sunset orange to cyan",
		FgStart:     "#FA6900",
		FgEnd:       "#69D2E7",
		BgColor:     "#FFFFFF",
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
	},
	{
		ID:          "emo",
		Name:        "Cheer Up 🎈",
		Description: "COLOURlovers classic: coral red to turquoise",
		FgStart:     "#FF6B6B",
		FgEnd:       "#4ECDC4",
		BgColor:     "#FFFFFF",
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
	},
	{
		ID:          "ocean",
		Name:        "Ocean Five 🌊",
		Description: "Deep sea teal to warm coral sun",
		FgStart:     "#00A0B0",
		FgEnd:       "#EB6841",
		BgColor:     "#FAF7F2",
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
	},
	{
		ID:          "thought",
		Name:        "Thought Provoking 🍷",
		Description: "Terracotta to ruby plum luxury",
		FgStart:     "#D95B43",
		FgEnd:       "#542437",
		BgColor:     "#FDF8F5",
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
	},
	{
		ID:          "cyber",
		Name:        "Cyber Neon ⚡",
		Description: "Futuristic neon cyan to hot pink on dark",
		FgStart:     "#00F5D4",
		FgEnd:       "#FF007F",
		BgColor:     "#0B0C10",
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
	},
	{
		ID:          "classic",
		Name:        "Classic Monochrome 🏁",
		Description: "High-contrast crisp black on white",
		FgStart:     "#000000",
		FgEnd:       "#000000",
		BgColor:     "#FFFFFF",
		Shape:       ShapeSquare,
		Gradient:    GradientDiagonal,
	},
}

type QRStyle struct {
	PresetIndex int
	Shape       ModuleShape
	Gradient    GradientDirection
	LogoPath    string
}

func defaultStyle() QRStyle {
	return QRStyle{
		PresetIndex: 0,
		Shape:       ShapeRounded,
		Gradient:    GradientDiagonal,
		LogoPath:    "",
	}
}

func (s QRStyle) CurrentPreset() StylePreset {
	if s.PresetIndex < 0 || s.PresetIndex >= len(presets) {
		return presets[0]
	}
	return presets[s.PresetIndex]
}

func parseHexColor(s string) (color.RGBA, error) {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	if len(s) != 6 {
		return color.RGBA{A: 255}, fmt.Errorf("invalid hex color %q", s)
	}
	val, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return color.RGBA{A: 255}, err
	}
	return color.RGBA{
		R: uint8(val >> 16),
		G: uint8((val >> 8) & 0xFF),
		B: uint8(val & 0xFF),
		A: 255,
	}, nil
}

func interpolateColor(c1, c2 color.RGBA, t float64) color.RGBA {
	if t <= 0 {
		return c1
	}
	if t >= 1 {
		return c2
	}
	r := float64(c1.R) + t*(float64(c2.R)-float64(c1.R))
	g := float64(c1.G) + t*(float64(c2.G)-float64(c1.G))
	b := float64(c1.B) + t*(float64(c2.B)-float64(c1.B))
	a := float64(c1.A) + t*(float64(c2.A)-float64(c1.A))
	return color.RGBA{
		R: uint8(r + 0.5),
		G: uint8(g + 0.5),
		B: uint8(b + 0.5),
		A: uint8(a + 0.5),
	}
}

func colorToHex(c color.RGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
