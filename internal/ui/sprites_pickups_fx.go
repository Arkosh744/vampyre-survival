package ui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// PickupSprite returns a cached sprite for the given pickup type and value.
// XP gems scale by value; HP pickups render as a pixel heart.
func PickupSprite(pickupType string, value int) *ebiten.Image {
	switch pickupType {
	case SpritePickupXP:
		if value >= 10 {
			return getCachedSprite("pickup_xp_large", 14, 14, drawPickupXPLarge)
		}
		return getCachedSprite("pickup_xp_small", 10, 10, drawPickupXPSmall)
	case SpritePickupHP:
		return getCachedSprite("pickup_hp", 16, 16, drawPickupHP)
	default:
		return getCachedSprite("pickup_xp_small", 10, 10, drawPickupXPSmall)
	}
}

func drawPickupXPLarge(img *ebiten.Image) {
	cx, cy := 7, 7
	base := ColorToRGBA(ColorCyan)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	for y := 0; y < 14; y++ {
		for x := 0; x < 14; x++ {
			dx := x - cx
			dy := y - cy
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			if dx+dy <= 5 {
				c := base
				if x < cx && y < cy {
					c = HighlightCyan
				} else if x > cx && y > cy {
					c = ShadowBlue
				}
				setPixelSafe(img, x, y, c)
			}
		}
	}

	// Sparkles at diamond tips
	setPixelSafe(img, 7, 1, white)
	setPixelSafe(img, 1, 7, white)
	setPixelSafe(img, 13, 7, white)
	setPixelSafe(img, 7, 13, white)
}

func drawPickupXPSmall(img *ebiten.Image) {
	cx, cy := 5, 5
	base := ColorToRGBA(ColorBlue)

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			dx := x - cx
			dy := y - cy
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			if dx+dy <= 3 {
				setPixelSafe(img, x, y, base)
			}
		}
	}

	// Single highlight pixel
	setPixelSafe(img, 4, 4, HighlightBlue)
}

func drawPickupHP(img *ebiten.Image) {
	base := ColorToRGBA(ColorRed)

	// Classic pixel heart shape, row by row
	heartRows := []struct {
		y         int
		colStart  int
		colEnd    int
		colStart2 int
		colEnd2   int
	}{
		{3, 3, 5, 10, 12},
		{4, 2, 6, 9, 13},
	}
	for _, r := range heartRows {
		for x := r.colStart; x <= r.colEnd; x++ {
			setPixelSafe(img, x, r.y, base)
		}
		for x := r.colStart2; x <= r.colEnd2; x++ {
			setPixelSafe(img, x, r.y, base)
		}
	}

	// Full-width rows
	fullRows := []struct {
		y        int
		colStart int
		colEnd   int
	}{
		{5, 1, 14},
		{6, 1, 14},
		{7, 1, 14},
		{8, 2, 13},
		{9, 3, 12},
		{10, 4, 11},
		{11, 5, 10},
		{12, 6, 9},
		{13, 7, 8},
	}
	for _, r := range fullRows {
		for x := r.colStart; x <= r.colEnd; x++ {
			setPixelSafe(img, x, r.y, base)
		}
	}

	// Outline
	drawOutline(img, ShadowRed)

	// Highlight on top-left lobe
	setPixelSafe(img, 3, 4, HighlightRed)
	setPixelSafe(img, 4, 4, HighlightRed)
	setPixelSafe(img, 3, 5, HighlightRed)
}

// ProjectileSprite returns a cached 8x8 projectile sprite with
// an oval shape, white-hot center, and orange shadow.
func ProjectileSprite() *ebiten.Image {
	return getCachedSprite("projectile", ProjectileSpriteSize, ProjectileSpriteSize, drawProjectile)
}

func drawProjectile(img *ebiten.Image) {
	base := ColorToRGBA(ColorYellow)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	orange := color.RGBA{R: 255, G: 160, B: 30, A: 255}

	// Oval shape: ((x-4)^2)/9 + ((y-4)^2)/4 <= 1
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			dx := float64(x) - 4.0
			dy := float64(y) - 4.0
			if (dx*dx)/9.0+(dy*dy)/4.0 <= 1.0 {
				c := base
				if x >= 5 {
					c = orange
				}
				setPixelSafe(img, x, y, c)
			}
		}
	}

	// White hot center
	setPixelSafe(img, 3, 3, white)
	setPixelSafe(img, 4, 3, white)
	setPixelSafe(img, 3, 4, white)
	setPixelSafe(img, 4, 4, white)
}

// WeaponEffectSprite returns a cached sprite for the given weapon
// effect kind: "lightning", "orbital", or a default star burst.
func WeaponEffectSprite(kind string) *ebiten.Image {
	switch kind {
	case "lightning":
		return getCachedSprite("fx_lightning", 10, 10, drawFXLightning)
	case "orbital":
		return getCachedSprite("fx_orbital", 12, 12, drawFXOrbital)
	default:
		return getCachedSprite("fx_default", 8, 8, drawFXDefault)
	}
}

func drawFXLightning(img *ebiten.Image) {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	cyan := ColorToRGBA(ColorCyan)
	glow := WithAlpha(HighlightCyan, 100)

	// Bolt pixel coordinates: (x, y, isWhite)
	type boltPx struct {
		x, y    int
		isWhite bool
	}
	bolt := []boltPx{
		{5, 0, true},
		{5, 1, true},
		{4, 2, false},
		{3, 3, false},
		{4, 4, true},
		{5, 5, true},
		{6, 6, false},
		{7, 7, false},
		{6, 8, false},
		{5, 9, false},
	}

	// Draw glow first (behind bolt pixels)
	for _, p := range bolt {
		setPixelSafe(img, p.x-1, p.y, glow)
		setPixelSafe(img, p.x+1, p.y, glow)
	}

	// Draw bolt pixels on top
	for _, p := range bolt {
		c := cyan
		if p.isWhite {
			c = white
		}
		setPixelSafe(img, p.x, p.y, c)
	}
}

func drawFXOrbital(img *ebiten.Image) {
	base := ColorToRGBA(ColorMagenta)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Purple filled circle
	drawCircleOnImage(img, 5.0, base)

	// Darker edge ring where 4.0 < dist <= 5.0
	cx, cy := 6.0, 6.0
	for py := 0; py < 12; py++ {
		for px := 0; px < 12; px++ {
			dx := float64(px) - cx
			dy := float64(py) - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > 4.0 && dist <= 5.0 {
				setPixelSafe(img, px, py, ShadowMagenta)
			}
		}
	}

	// White center glow
	drawRectOn(img, 5, 5, 2, 2, white)
}

func drawFXDefault(img *ebiten.Image) {
	yellow := ColorToRGBA(ColorYellow)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Cross shape: horizontal and vertical arms
	drawRectOn(img, 0, 3, 8, 2, yellow)
	drawRectOn(img, 3, 0, 2, 8, yellow)

	// White center 2x2
	drawRectOn(img, 3, 3, 2, 2, white)
}

// DeathParticleSprite returns a cached 6x6 star-burst particle
// tinted with the given color, used for enemy death effects.
func DeathParticleSprite(clr color.RGBA) *ebiten.Image {
	key := fmt.Sprintf("death_%d_%d_%d", clr.R, clr.G, clr.B)
	return getCachedSprite(key, 6, 6, func(img *ebiten.Image) {
		drawDeathParticle(img, clr)
	})
}

func drawDeathParticle(img *ebiten.Image, clr color.RGBA) {
	shadow := ShadowOf(clr)
	highlight := HighlightOf(clr)

	// Center 2x2
	setPixelSafe(img, 2, 2, clr)
	setPixelSafe(img, 3, 2, clr)
	setPixelSafe(img, 2, 3, clr)
	setPixelSafe(img, 3, 3, clr)

	// Horizontal arms (shadow)
	setPixelSafe(img, 0, 2, shadow)
	setPixelSafe(img, 1, 2, shadow)
	setPixelSafe(img, 4, 2, shadow)
	setPixelSafe(img, 5, 2, shadow)
	setPixelSafe(img, 0, 3, shadow)
	setPixelSafe(img, 1, 3, shadow)
	setPixelSafe(img, 4, 3, shadow)
	setPixelSafe(img, 5, 3, shadow)

	// Vertical arms (shadow)
	setPixelSafe(img, 2, 0, shadow)
	setPixelSafe(img, 3, 0, shadow)
	setPixelSafe(img, 2, 1, shadow)
	setPixelSafe(img, 3, 1, shadow)
	setPixelSafe(img, 2, 4, shadow)
	setPixelSafe(img, 3, 4, shadow)
	setPixelSafe(img, 2, 5, shadow)
	setPixelSafe(img, 3, 5, shadow)

	// Bright center highlights
	setPixelSafe(img, 2, 2, highlight)
	setPixelSafe(img, 3, 3, highlight)
}
