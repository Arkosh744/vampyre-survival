package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Grass tile colors.
var (
	grassBase      = color.RGBA{R: 25, G: 55, B: 25, A: 255}
	grassDark      = color.RGBA{R: 30, G: 70, B: 30, A: 255}
	grassMid       = color.RGBA{R: 40, G: 100, B: 35, A: 255}
	grassLight     = color.RGBA{R: 55, G: 130, B: 45, A: 255}
	grassYellow    = color.RGBA{R: 90, G: 120, B: 40, A: 255}
	grassDarkBase  = color.RGBA{R: 20, G: 48, B: 20, A: 255}
	grassLightBase = color.RGBA{R: 30, G: 62, B: 28, A: 255}
)

// Blade describes a single grass blade for deterministic tile generation.
type blade struct {
	x, y, h int
	clr      color.RGBA
	tipClr   color.RGBA
	width    int
}

// Grass variant definitions: each variant has a set of blades.
var grassVariants = [4][]blade{
	{ // Variant 0: 4 blades, medium density
		{x: 3, y: 12, h: 5, clr: grassDark, tipClr: grassLight, width: 1},
		{x: 7, y: 10, h: 7, clr: grassMid, tipClr: grassLight, width: 1},
		{x: 11, y: 11, h: 4, clr: grassDark, tipClr: grassMid, width: 1},
		{x: 14, y: 13, h: 3, clr: grassMid, tipClr: grassYellow, width: 1},
	},
	{ // Variant 1: 5 blades, dense
		{x: 2, y: 10, h: 6, clr: grassMid, tipClr: grassLight, width: 1},
		{x: 5, y: 11, h: 5, clr: grassDark, tipClr: grassMid, width: 2},
		{x: 8, y: 9, h: 7, clr: grassMid, tipClr: grassYellow, width: 1},
		{x: 12, y: 12, h: 4, clr: grassDark, tipClr: grassLight, width: 1},
		{x: 14, y: 10, h: 5, clr: grassMid, tipClr: grassLight, width: 1},
	},
	{ // Variant 2: 3 blades, sparse
		{x: 4, y: 11, h: 5, clr: grassDark, tipClr: grassMid, width: 1},
		{x: 9, y: 10, h: 6, clr: grassMid, tipClr: grassLight, width: 2},
		{x: 13, y: 12, h: 3, clr: grassDark, tipClr: grassYellow, width: 1},
	},
	{ // Variant 3: 5 blades, mixed with yellow tips
		{x: 1, y: 11, h: 5, clr: grassMid, tipClr: grassYellow, width: 1},
		{x: 4, y: 9, h: 7, clr: grassDark, tipClr: grassLight, width: 1},
		{x: 7, y: 12, h: 4, clr: grassMid, tipClr: grassMid, width: 2},
		{x: 10, y: 10, h: 6, clr: grassDark, tipClr: grassLight, width: 1},
		{x: 14, y: 11, h: 4, clr: grassMid, tipClr: grassYellow, width: 1},
	},
}

// GrassSprite returns a cached 16x16 grass tile for the given variant (0-3).
func GrassSprite(variant int) *ebiten.Image {
	v := variant % 4
	key := "grass_" + string(rune('0'+v))
	return getCachedSprite(key, 16, 16, func(img *ebiten.Image) {
		// Fill base
		base := grassBase
		if v%2 == 1 {
			base = grassDarkBase
		}
		img.Fill(base)

		// Subtle ground variation: a few scattered darker/lighter pixels
		if v == 0 || v == 3 {
			setPixelSafe(img, 6, 14, grassLightBase)
			setPixelSafe(img, 10, 7, grassLightBase)
		}
		if v == 1 || v == 2 {
			setPixelSafe(img, 3, 5, grassDarkBase)
			setPixelSafe(img, 13, 8, grassLightBase)
		}

		// Draw grass blades
		for _, b := range grassVariants[v] {
			for dy := 0; dy < b.h; dy++ {
				py := b.y - dy
				c := b.clr
				if dy >= b.h-1 {
					c = b.tipClr
				}
				for dx := 0; dx < b.width; dx++ {
					setPixelSafe(img, b.x+dx, py, c)
				}
			}
		}
	})
}
