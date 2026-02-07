package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// BossEnemySprite returns a cached pixel-art sprite for boss-tier enemies.
// Returns nil for unrecognized enemy types.
func BossEnemySprite(enemyType string, variant bool) *ebiten.Image {
	switch enemyType {

	case SpriteEnemyBoss:
		return getCachedSprite("enemy_boss", BossSpriteSize, BossSpriteSize, func(img *ebiten.Image) {
			// Cape (rows 18-44): dark red flowing cape
			drawCape(img, 8, 18, 32, 26, color.RGBA{140, 20, 20, 255}, ShadowBoldRed)

			// Body/torso (rows 14-30): dark garment
			drawRectOn(img, 16, 14, 16, 16, color.RGBA{60, 20, 20, 255})

			// Head (rows 6-14): pale face
			drawRectOn(img, 18, 6, 12, 8, ColorSkin)
			// Right shadow on face
			drawRectOn(img, 28, 6, 2, 8, ColorSkinShadow)

			// Hair (rows 4-7): dark hair
			drawRectOn(img, 17, 4, 14, 3, color.RGBA{30, 20, 20, 255})

			// Eyes (row 9): red glowing
			drawEyes(img, 24, 9, 3, color.RGBA{255, 40, 40, 255}, color.RGBA{180, 0, 0, 255})

			// Fangs: white pixels
			setPixelSafe(img, 22, 13, color.RGBA{255, 255, 255, 255})
			setPixelSafe(img, 25, 13, color.RGBA{255, 255, 255, 255})

			// Crown/horns: golden spikes
			hornClr := color.RGBA{255, 220, 50, 255}
			// Left horn
			setPixelSafe(img, 20, 3, hornClr)
			setPixelSafe(img, 20, 2, hornClr)
			// Center horn (tallest)
			setPixelSafe(img, 24, 2, hornClr)
			setPixelSafe(img, 24, 1, hornClr)
			// Right horn
			setPixelSafe(img, 28, 3, hornClr)
			setPixelSafe(img, 28, 2, hornClr)

			// Hands: skin-colored 2x3 blocks on each side
			drawRectOn(img, 14, 22, 2, 3, ColorSkin)
			drawRectOn(img, 32, 22, 2, 3, ColorSkin)

			// Feet: dark boots
			bootClr := color.RGBA{40, 20, 20, 255}
			drawRectOn(img, 18, 42, 4, 2, bootClr)
			drawRectOn(img, 27, 42, 4, 2, bootClr)

			// Outline
			drawOutline(img, ColorOutline)
		})

	case SpriteEnemyElder:
		return getCachedSprite("enemy_elder", ElderSpriteSize, ElderSpriteSize, func(img *ebiten.Image) {
			drawElderSprite(img, elderPalette{
				capeBase:    color.RGBA{80, 20, 100, 255},
				capeShadow:  ShadowMagenta,
				vPatternClr: color.RGBA{200, 50, 50, 255},
				armorBase:   color.RGBA{70, 70, 80, 255},
				armorHL:     color.RGBA{120, 120, 140, 255},
				eyeClr:      color.RGBA{255, 40, 40, 255},
				pupilClr:    color.RGBA{180, 0, 0, 255},
				crownClr:    color.RGBA{255, 220, 50, 255},
				outlineClr:  ColorOutline,
			})
		})

	case SpriteEnemyElderEnraged:
		return getCachedSprite("enemy_elder_enraged", ElderSpriteSize, ElderSpriteSize, func(img *ebiten.Image) {
			drawElderSprite(img, elderPalette{
				capeBase:    color.RGBA{120, 20, 30, 255},
				capeShadow:  ShadowBoldRed,
				vPatternClr: color.RGBA{255, 180, 50, 255},
				armorBase:   color.RGBA{90, 50, 50, 255},
				armorHL:     color.RGBA{120, 120, 140, 255},
				eyeClr:      color.RGBA{255, 100, 100, 255},
				pupilClr:    color.RGBA{255, 200, 200, 255},
				crownClr:    color.RGBA{255, 240, 100, 255},
				outlineClr:  color.RGBA{80, 20, 20, 255},
			})

			// Red glow overlay: shift R channel up by 30 on all visible pixels
			b := img.Bounds()
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					r, g, bb, a := img.At(x, y).RGBA()
					if a > 0 {
						nr := uint8(min(int(r>>8)+30, 255))
						img.Set(x, y, color.RGBA{R: nr, G: uint8(g >> 8), B: uint8(bb >> 8), A: uint8(a >> 8)})
					}
				}
			}
		})

	default:
		return nil
	}
}

// elderPalette holds configurable colors for the Elder sprite variants.
type elderPalette struct {
	capeBase    color.RGBA
	capeShadow  color.RGBA
	vPatternClr color.RGBA
	armorBase   color.RGBA
	armorHL     color.RGBA
	eyeClr      color.RGBA
	pupilClr    color.RGBA
	crownClr    color.RGBA
	outlineClr  color.RGBA
}

// drawElderSprite draws the full Elder (Vampire Lord) sprite onto img
// using the provided palette.
func drawElderSprite(img *ebiten.Image, p elderPalette) {
	// Cape (rows 20-60): dark purple, spanning almost full width
	drawCape(img, 4, 20, 56, 40, p.capeBase, p.capeShadow)

	// V-pattern ornament on cape (rows 30-55)
	for row := 30; row <= 55; row++ {
		t := float64(row-30) / 25.0
		leftCol := 20 + int(t*12)
		rightCol := 44 - int(t*12)
		// 2px wide V arms
		setPixelSafe(img, leftCol, row, p.vPatternClr)
		setPixelSafe(img, leftCol+1, row, p.vPatternClr)
		setPixelSafe(img, rightCol, row, p.vPatternClr)
		setPixelSafe(img, rightCol+1, row, p.vPatternClr)
	}

	// Armored chest (rows 18-30): dark gray armor
	drawRectOn(img, 20, 18, 24, 12, p.armorBase)
	// Highlight edges: top row
	for x := 20; x < 44; x++ {
		setPixelSafe(img, x, 18, p.armorHL)
	}
	// Highlight edges: left col
	for y := 18; y < 30; y++ {
		setPixelSafe(img, 20, y, p.armorHL)
	}

	// Head (rows 8-18): pale face
	drawRectOn(img, 24, 8, 16, 10, ColorSkin)
	// Right shadow
	drawRectOn(img, 38, 8, 2, 10, ColorSkinShadow)

	// Eyes (row 12): glowing red
	drawEyes(img, 32, 12, 4, p.eyeClr, p.pupilClr)

	// Fangs: white pixels
	setPixelSafe(img, 30, 17, color.RGBA{255, 255, 255, 255})
	setPixelSafe(img, 33, 17, color.RGBA{255, 255, 255, 255})

	// Crown: base bar
	drawRectOn(img, 23, 7, 18, 1, p.crownClr)
	// Crown: 5 golden points at x=24,27,32,37,40
	// Side spikes (shortest)
	setPixelSafe(img, 24, 7, p.crownClr)
	setPixelSafe(img, 24, 6, p.crownClr)
	setPixelSafe(img, 40, 7, p.crownClr)
	setPixelSafe(img, 40, 6, p.crownClr)
	// Inner spikes (medium)
	setPixelSafe(img, 27, 7, p.crownClr)
	setPixelSafe(img, 27, 6, p.crownClr)
	setPixelSafe(img, 27, 5, p.crownClr)
	setPixelSafe(img, 37, 7, p.crownClr)
	setPixelSafe(img, 37, 6, p.crownClr)
	setPixelSafe(img, 37, 5, p.crownClr)
	// Center spike (tallest)
	setPixelSafe(img, 32, 7, p.crownClr)
	setPixelSafe(img, 32, 6, p.crownClr)
	setPixelSafe(img, 32, 5, p.crownClr)
	setPixelSafe(img, 32, 4, p.crownClr)

	// Clawed hands: skin-colored 4x5 blocks
	drawRectOn(img, 10, 28, 4, 5, ColorSkin)
	drawRectOn(img, 50, 28, 4, 5, ColorSkin)
	// Claws: darker tips below hands
	clawClr := color.RGBA{200, 180, 160, 255}
	setPixelSafe(img, 10, 33, clawClr)
	setPixelSafe(img, 12, 33, clawClr)
	setPixelSafe(img, 50, 33, clawClr)
	setPixelSafe(img, 52, 33, clawClr)

	// Outline
	drawOutline(img, p.outlineClr)
}
