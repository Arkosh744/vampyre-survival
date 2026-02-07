package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// EnemySprite returns a cached pixel-art sprite for the given enemy type.
// Boss, elder, and elder_enraged sprites are handled in a separate file.
// The variant flag activates an alternate visual (e.g. dashing state).
func EnemySprite(enemyType string, variant bool) *ebiten.Image {
	switch enemyType {

	case SpriteEnemyNormal:
		return getCachedSprite("enemy_normal", NormalSpriteSize, NormalSpriteSize, func(img *ebiten.Image) {
			baseClr := ColorToRGBA(ColorRed)

			// Body: red oval 8x8 in center (rows 8-15, cols 8-15)
			for y := 8; y < 16; y++ {
				for x := 8; x < 16; x++ {
					c := baseClr
					if y >= 14 {
						c = ShadowRed
					}
					setPixelSafe(img, x, y, c)
				}
			}

			// Wings
			drawWings(img, 12, 11, 7, baseClr, ShadowRed)

			// Eyes: red glow
			drawEyes(img, 12, 9, 2, color.RGBA{255, 80, 80, 255}, color.RGBA{200, 0, 0, 255})

			// Fangs: 2 white pixels below body center
			setPixelSafe(img, 11, 16, color.RGBA{255, 255, 255, 255})
			setPixelSafe(img, 12, 16, color.RGBA{255, 255, 255, 255})

			// Shadow below body
			drawRectOn(img, 9, 20, 6, 2, color.RGBA{30, 10, 10, 120})

			// Outline
			drawOutline(img, ColorOutline)
		})

	case SpriteEnemySwarmer:
		return getCachedSprite("enemy_swarmer", SwarmerSpriteSize, SwarmerSpriteSize, func(img *ebiten.Image) {
			baseClr := ColorToRGBA(ColorYellow)

			// Body: 4x4 center (rows 4-7, cols 4-7)
			for y := 4; y < 8; y++ {
				for x := 4; x < 8; x++ {
					setPixelSafe(img, x, y, baseClr)
				}
			}

			// Wings
			drawWings(img, 6, 5, 3, baseClr, ShadowYellow)

			// Eyes: 2 red pixels
			setPixelSafe(img, 5, 4, color.RGBA{255, 50, 50, 255})
			setPixelSafe(img, 6, 4, color.RGBA{255, 50, 50, 255})
		})

	case SpriteEnemyTank:
		return getCachedSprite("enemy_tank", TankSpriteSize, TankSpriteSize, func(img *ebiten.Image) {
			baseClr := ColorToRGBA(ColorMagenta)

			// Body: 20x24 block (rows 4-27, cols 8-27) with checkerboard dither
			for y := 4; y < 28; y++ {
				for x := 8; x < 28; x++ {
					c := baseClr
					if (x+y)%2 == 0 {
						c = ShadowMagenta
					}
					setPixelSafe(img, x, y, c)
				}
			}

			// Head: 8x6 on top (rows 2-7, cols 14-21)
			for y := 2; y < 8; y++ {
				for x := 14; x < 22; x++ {
					c := baseClr
					if (x+y)%2 == 0 {
						c = ShadowMagenta
					}
					setPixelSafe(img, x, y, c)
				}
			}

			// Shoulders: extra 2px on each side of body top (rows 4-6)
			for y := 4; y < 7; y++ {
				for x := 6; x < 8; x++ {
					c := baseClr
					if (x+y)%2 == 0 {
						c = ShadowMagenta
					}
					setPixelSafe(img, x, y, c)
				}
				for x := 28; x < 30; x++ {
					c := baseClr
					if (x+y)%2 == 0 {
						c = ShadowMagenta
					}
					setPixelSafe(img, x, y, c)
				}
			}

			// Legs: 2 thick legs, 4px wide each (rows 28-33)
			for y := 28; y < 34; y++ {
				for x := 11; x < 15; x++ {
					setPixelSafe(img, x, y, ShadowMagenta)
				}
				for x := 22; x < 26; x++ {
					setPixelSafe(img, x, y, ShadowMagenta)
				}
			}

			// Eyes: glowing cyan
			drawEyes(img, 18, 4, 2, HighlightCyan, color.RGBA{0, 200, 200, 255})

			// Armor line: horizontal dark stripe at row 16
			for x := 8; x < 28; x++ {
				setPixelSafe(img, x, 16, color.RGBA{80, 30, 80, 255})
			}

			// Outline
			drawOutline(img, ColorOutline)
		})

	case SpriteEnemyDasher:
		if variant {
			return getCachedSprite("enemy_dasher_active", DasherSpriteSize, DasherSpriteSize, func(img *ebiten.Image) {
				// Dashing variant: brighter with motion blur

				// Spectral shape using HighlightCyan as base everywhere
				for y := 2; y < 6; y++ {
					drawRectOn(img, 4, y, 12, 1, HighlightCyan)
				}
				for y := 6; y < 10; y++ {
					drawRectOn(img, 5, y, 10, 1, HighlightCyan)
				}
				for y := 10; y < 14; y++ {
					drawRectOn(img, 6, y, 8, 1, HighlightCyan)
				}
				for y := 14; y < 17; y++ {
					drawRectOn(img, 7, y, 6, 1, WithAlpha(HighlightCyan, 180))
				}
				for y := 17; y < 19; y++ {
					drawRectOn(img, 8, y, 4, 1, WithAlpha(HighlightCyan, 100))
				}
				// Trailing wisps
				setPixelSafe(img, 8, 19, WithAlpha(HighlightCyan, 60))
				setPixelSafe(img, 11, 19, WithAlpha(HighlightCyan, 60))

				// Motion blur streaks at rows 6, 10, 14
				motionClr := WithAlpha(HighlightCyan, 120)
				for x := 2; x < 5; x++ {
					setPixelSafe(img, x, 6, motionClr)
					setPixelSafe(img, x, 10, motionClr)
					setPixelSafe(img, x, 14, motionClr)
				}

				// Eyes: bright white
				setPixelSafe(img, 8, 5, color.RGBA{255, 255, 255, 255})
				setPixelSafe(img, 11, 5, color.RGBA{255, 255, 255, 255})
			})
		}

		return getCachedSprite("enemy_dasher", DasherSpriteSize, DasherSpriteSize, func(img *ebiten.Image) {
			cyanBase := ColorToRGBA(ColorCyan)

			// Spectral shape: wider at top, tapers to bottom
			for y := 2; y < 6; y++ {
				drawRectOn(img, 4, y, 12, 1, HighlightCyan)
			}
			for y := 6; y < 10; y++ {
				drawRectOn(img, 5, y, 10, 1, cyanBase)
			}
			for y := 10; y < 14; y++ {
				drawRectOn(img, 6, y, 8, 1, ShadowCyan)
			}
			for y := 14; y < 17; y++ {
				drawRectOn(img, 7, y, 6, 1, WithAlpha(cyanBase, 180))
			}
			for y := 17; y < 19; y++ {
				drawRectOn(img, 8, y, 4, 1, WithAlpha(cyanBase, 100))
			}
			// Trailing wisps
			setPixelSafe(img, 8, 19, WithAlpha(cyanBase, 60))
			setPixelSafe(img, 11, 19, WithAlpha(cyanBase, 60))

			// Eyes: bright white
			setPixelSafe(img, 8, 5, color.RGBA{255, 255, 255, 255})
			setPixelSafe(img, 11, 5, color.RGBA{255, 255, 255, 255})
		})

	case SpriteEnemyBoss, SpriteEnemyElder, SpriteEnemyElderEnraged:
		return BossEnemySprite(enemyType, variant)

	default:
		// Normal enemy as fallback
		return getCachedSprite("enemy_normal", NormalSpriteSize, NormalSpriteSize, func(img *ebiten.Image) {
			baseClr := ColorToRGBA(ColorRed)
			for y := 8; y < 16; y++ {
				for x := 8; x < 16; x++ {
					c := baseClr
					if y >= 14 {
						c = ShadowRed
					}
					setPixelSafe(img, x, y, c)
				}
			}
			drawWings(img, 12, 11, 7, baseClr, ShadowRed)
			drawEyes(img, 12, 9, 2, color.RGBA{255, 80, 80, 255}, color.RGBA{200, 0, 0, 255})
			setPixelSafe(img, 11, 16, color.RGBA{255, 255, 255, 255})
			setPixelSafe(img, 12, 16, color.RGBA{255, 255, 255, 255})
			drawRectOn(img, 9, 20, 6, 2, color.RGBA{30, 10, 10, 120})
			drawOutline(img, ColorOutline)
		})
	}
}
