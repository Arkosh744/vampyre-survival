package ui

import "image/color"

// Color constants used as string keys for the color palette.
// Originally ANSI escape codes; now simply palette identifiers
// mapped to RGBA by ColorToRGBA.
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBoldRed = "\033[1;31m"
	ColorDim     = "\033[2m"
)

// Shading triplets: base, shadow, highlight for each palette color.
var (
	ShadowRed     = color.RGBA{R: 150, G: 30, B: 30, A: 255}
	HighlightRed  = color.RGBA{R: 255, G: 100, B: 90, A: 255}

	ShadowGreen     = color.RGBA{R: 20, G: 140, B: 20, A: 255}
	HighlightGreen  = color.RGBA{R: 100, G: 255, B: 100, A: 255}

	ShadowYellow     = color.RGBA{R: 200, G: 170, B: 30, A: 255}
	HighlightYellow  = color.RGBA{R: 255, G: 245, B: 150, A: 255}

	ShadowBlue     = color.RGBA{R: 30, G: 60, B: 150, A: 255}
	HighlightBlue  = color.RGBA{R: 100, G: 150, B: 255, A: 255}

	ShadowMagenta     = color.RGBA{R: 140, G: 30, B: 140, A: 255}
	HighlightMagenta  = color.RGBA{R: 240, G: 120, B: 240, A: 255}

	ShadowCyan     = color.RGBA{R: 20, G: 150, B: 150, A: 255}
	HighlightCyan  = color.RGBA{R: 100, G: 255, B: 255, A: 255}

	ShadowBoldRed     = color.RGBA{R: 180, G: 30, B: 30, A: 255}
	HighlightBoldRed  = color.RGBA{R: 255, G: 130, B: 130, A: 255}
)

// Skin and outline colors for humanoid sprites.
var (
	ColorSkin       = color.RGBA{R: 228, G: 190, B: 160, A: 255}
	ColorSkinShadow = color.RGBA{R: 190, G: 150, B: 120, A: 255}
	ColorOutline    = color.RGBA{R: 42, G: 42, B: 42, A: 255}
)

// ShadowOf darkens a color by ~35%.
func ShadowOf(c color.RGBA) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c.R) * 0.65),
		G: uint8(float64(c.G) * 0.65),
		B: uint8(float64(c.B) * 0.65),
		A: c.A,
	}
}

// HighlightOf lightens a color by ~40%.
func HighlightOf(c color.RGBA) color.RGBA {
	r := float64(c.R) + (255-float64(c.R))*0.4
	g := float64(c.G) + (255-float64(c.G))*0.4
	b := float64(c.B) + (255-float64(c.B))*0.4
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: c.A,
	}
}

// WithAlpha returns a copy of the color with the given alpha value.
func WithAlpha(c color.RGBA, a uint8) color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: a}
}

// ColorToRGBA maps color string constants to RGBA values
// suitable for Ebitengine rendering.
func ColorToRGBA(ansiColor string) color.RGBA {
	switch ansiColor {
	case ColorRed:
		return color.RGBA{R: 220, G: 50, B: 47, A: 255}
	case ColorGreen:
		return color.RGBA{R: 38, G: 210, B: 38, A: 255}
	case ColorYellow:
		return color.RGBA{R: 255, G: 220, B: 50, A: 255}
	case ColorBlue:
		return color.RGBA{R: 50, G: 100, B: 220, A: 255}
	case ColorMagenta:
		return color.RGBA{R: 200, G: 60, B: 200, A: 255}
	case ColorCyan:
		return color.RGBA{R: 42, G: 220, B: 220, A: 255}
	case ColorBoldRed:
		return color.RGBA{R: 255, G: 60, B: 60, A: 255}
	case ColorDim:
		return color.RGBA{R: 100, G: 100, B: 100, A: 255}
	case ColorReset, ColorWhite:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255}
	default:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}
