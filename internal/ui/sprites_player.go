package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// PlayerSprite returns a cached 24x24 vampire hunter sprite.
// When invulnerable is true, the sprite gets a golden aura and glowing eyes.
func PlayerSprite(invulnerable bool) *ebiten.Image {
	if invulnerable {
		return getCachedSprite("player_invuln", PlayerSpriteSize, PlayerSpriteSize, drawPlayerInvulnerable)
	}
	return getCachedSprite("player", PlayerSpriteSize, PlayerSpriteSize, drawPlayerNormal)
}

// PlayerDirectionIndicator returns a cached 6x6 chevron pointing right.
// Normal variant is green; invulnerable variant is yellow.
func PlayerDirectionIndicator(invulnerable bool) *ebiten.Image {
	if invulnerable {
		return getCachedSprite("player_dir_invuln", DirIndicatorSize, DirIndicatorSize, drawDirIndicatorInvuln)
	}
	return getCachedSprite("player_dir", DirIndicatorSize, DirIndicatorSize, drawDirIndicatorNormal)
}

func drawPlayerBase(img *ebiten.Image, eyeClr color.RGBA) {
	hairClr := color.RGBA{R: 60, G: 40, B: 30, A: 255}
	bodyClr := color.RGBA{R: 30, G: 80, B: 80, A: 255}
	bodyShadowClr := color.RGBA{R: 20, G: 55, B: 55, A: 255}
	capeClr := color.RGBA{R: 50, G: 30, B: 60, A: 255}
	capeShadowClr := color.RGBA{R: 30, G: 15, B: 35, A: 255}
	bootClr := color.RGBA{R: 60, G: 40, B: 30, A: 255}

	// Hair (rows 2-4, ~8px wide centered)
	drawRectOn(img, 8, 2, 8, 2, hairClr)

	// Head (rows 4-8, 6px wide centered at col 9-14)
	drawRectOn(img, 9, 4, 6, 4, ColorSkin)
	// Skin shadow on right edge
	for y := 4; y < 8; y++ {
		setPixelSafe(img, 14, y, ColorSkinShadow)
	}

	// Eyes (row 6)
	drawEyes(img, 12, 6, 3, eyeClr, ColorOutline)

	// Body/torso (rows 8-16, 10px wide at col 7-16)
	drawRectOn(img, 7, 8, 10, 8, bodyClr)
	// Shadow on left edge (col 7) and right edge (col 16)
	for y := 8; y < 16; y++ {
		setPixelSafe(img, 7, y, bodyShadowClr)
		setPixelSafe(img, 16, y, bodyShadowClr)
	}

	// Cape (rows 10-22, extends wider than body)
	drawCape(img, 5, 10, 14, 12, capeClr, capeShadowClr)

	// Arms (rows 10-14, 2px wide on each side)
	drawRectOn(img, 6, 10, 2, 4, ColorSkin)
	drawRectOn(img, 17, 10, 2, 4, ColorSkin)

	// Feet/boots (rows 20-22)
	drawRectOn(img, 8, 20, 3, 2, bootClr)
	drawRectOn(img, 14, 20, 3, 2, bootClr)
}

func drawPlayerNormal(img *ebiten.Image) {
	whiteEye := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	drawPlayerBase(img, whiteEye)
	drawOutline(img, ColorOutline)
}

func drawPlayerInvulnerable(img *ebiten.Image) {
	goldenEye := color.RGBA{R: 255, G: 220, B: 50, A: 255}
	goldenOutline := color.RGBA{R: 255, G: 220, B: 50, A: 200}
	drawPlayerBase(img, goldenEye)
	drawOutline(img, goldenOutline)
}

func drawDirChevron(img *ebiten.Image, baseClr, highlightClr color.RGBA) {
	// Row 0: col 2
	setPixelSafe(img, 2, 0, baseClr)
	// Row 1: col 3, 4
	setPixelSafe(img, 3, 1, baseClr)
	setPixelSafe(img, 4, 1, highlightClr)
	// Row 2: col 4, 5
	setPixelSafe(img, 4, 2, highlightClr)
	setPixelSafe(img, 5, 2, baseClr)
	// Row 3: col 3, 4
	setPixelSafe(img, 3, 3, baseClr)
	setPixelSafe(img, 4, 3, highlightClr)
	// Row 4: col 2
	setPixelSafe(img, 2, 4, baseClr)
}

func drawDirIndicatorNormal(img *ebiten.Image) {
	drawDirChevron(img, ColorToRGBA(ColorGreen), HighlightGreen)
}

func drawDirIndicatorInvuln(img *ebiten.Image) {
	drawDirChevron(img, ColorToRGBA(ColorYellow), HighlightYellow)
}
