package ui

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite sizes in pixels for each entity type.
const (
	PlayerSpriteSize     = 24
	BossSpriteSize       = 48
	SwarmerSpriteSize    = 12
	TankSpriteSize       = 36
	DasherSpriteSize     = 20
	ElderSpriteSize      = 64
	NormalSpriteSize     = 24
	ProjectileSpriteSize = 8
	DirIndicatorSize     = 6
)

// Enemy type string constants (mirrors entity package to avoid import cycle).
const (
	SpriteEnemyNormal       = "normal"
	SpriteEnemyBoss         = "boss"
	SpriteEnemySwarmer      = "swarmer"
	SpriteEnemyTank         = "tank"
	SpriteEnemyDasher       = "dasher"
	SpriteEnemyElder        = "elder"
	SpriteEnemyElderEnraged = "elder_enraged"
)

// Pickup type string constants.
const (
	SpritePickupXP = "xp"
	SpritePickupHP = "hp"
)

var (
	spriteCache   = make(map[string]*ebiten.Image)
	spriteCacheMu sync.Mutex
)

// getCachedSprite returns a cached sprite by key, creating it on first access
// via the provided draw function.
func getCachedSprite(key string, w, h int, draw func(img *ebiten.Image)) *ebiten.Image {
	spriteCacheMu.Lock()
	defer spriteCacheMu.Unlock()
	if img, ok := spriteCache[key]; ok {
		return img
	}
	img := ebiten.NewImage(w, h)
	draw(img)
	spriteCache[key] = img
	return img
}

// --- Pixel-art helper functions ---

// setPixelSafe sets a pixel only if (x,y) is within image bounds.
func setPixelSafe(img *ebiten.Image, x, y int, clr color.RGBA) {
	b := img.Bounds()
	if x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y {
		img.Set(x, y, clr)
	}
}

// drawRectOn fills a sub-rectangle on an existing image.
func drawRectOn(img *ebiten.Image, x, y, w, h int, clr color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			setPixelSafe(img, px, py, clr)
		}
	}
}

// drawOutline draws a 1px dark border around all non-transparent pixels.
func drawOutline(img *ebiten.Image, clr color.RGBA) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	// Collect non-transparent pixel positions
	filled := make([][]bool, h)
	for y := 0; y < h; y++ {
		filled[y] = make([]bool, w)
		for x := 0; x < w; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			filled[y][x] = a > 0
		}
	}
	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if filled[y][x] {
				continue
			}
			for d := 0; d < 4; d++ {
				nx, ny := x+dx[d], y+dy[d]
				if nx >= 0 && nx < w && ny >= 0 && ny < h && filled[ny][nx] {
					img.Set(x, y, clr)
					break
				}
			}
		}
	}
}

// drawEyes draws a pair of 2px eyes with 1px pupils.
func drawEyes(img *ebiten.Image, cx, cy, spacing int, eyeClr, pupilClr color.RGBA) {
	lx := cx - spacing
	rx := cx + spacing - 1
	// Left eye 2x2
	setPixelSafe(img, lx, cy, eyeClr)
	setPixelSafe(img, lx+1, cy, eyeClr)
	setPixelSafe(img, lx, cy+1, eyeClr)
	setPixelSafe(img, lx+1, cy+1, eyeClr)
	// Left pupil
	setPixelSafe(img, lx+1, cy+1, pupilClr)
	// Right eye 2x2
	setPixelSafe(img, rx, cy, eyeClr)
	setPixelSafe(img, rx+1, cy, eyeClr)
	setPixelSafe(img, rx, cy+1, eyeClr)
	setPixelSafe(img, rx+1, cy+1, eyeClr)
	// Right pupil
	setPixelSafe(img, rx, cy+1, pupilClr)
}

// drawWings draws bat-like wings extending from center.
func drawWings(img *ebiten.Image, cx, cy, span int, clr, shadowClr color.RGBA) {
	for i := 1; i <= span; i++ {
		// Wing height tapers toward tips
		wingH := span - i/2
		if wingH < 1 {
			wingH = 1
		}
		for dy := 0; dy < wingH; dy++ {
			c := clr
			if dy >= wingH-1 {
				c = shadowClr
			}
			// Left wing
			setPixelSafe(img, cx-i, cy-wingH/2+dy, c)
			// Right wing
			setPixelSafe(img, cx+i-1, cy-wingH/2+dy, c)
		}
	}
}

// drawCape draws a flowing cape shape.
func drawCape(img *ebiten.Image, x, y, w, h int, clr, shadowClr color.RGBA) {
	for row := 0; row < h; row++ {
		// Cape widens slightly toward bottom
		extra := row / 4
		for col := -extra; col < w+extra; col++ {
			px := x + col
			c := clr
			// Shadow on edges and bottom
			if col < 1 || col >= w-1 || row >= h-2 {
				c = shadowClr
			}
			setPixelSafe(img, px, y+row, c)
		}
	}
}

// --- Legacy shape helpers ---

// drawCircleOnImage fills a circle centered in the given image.
func drawCircleOnImage(img *ebiten.Image, radius float64, clr color.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	cx := float64(w) / 2
	cy := float64(h) / 2
	r2 := radius * radius
	for py := 0; py < h; py++ {
		dy := float64(py) - cy
		for px := 0; px < w; px++ {
			dx := float64(px) - cx
			if dx*dx+dy*dy <= r2 {
				img.Set(px, py, clr)
			}
		}
	}
}

// drawCrossOnImage draws a plus/cross shape centered in the given image.
func drawCrossOnImage(img *ebiten.Image, size int, clr color.RGBA) {
	thickness := size / 3
	if thickness < 2 {
		thickness = 2
	}
	cx := size / 2
	cy := size / 2

	for py := cy - thickness/2; py < cy+thickness/2; py++ {
		for px := 0; px < size; px++ {
			if px >= 0 && px < size && py >= 0 && py < size {
				img.Set(px, py, clr)
			}
		}
	}
	for py := 0; py < size; py++ {
		for px := cx - thickness/2; px < cx+thickness/2; px++ {
			if px >= 0 && px < size && py >= 0 && py < size {
				img.Set(px, py, clr)
			}
		}
	}
}

// drawVPatternOnImage draws a V-shaped pattern centered in the given image.
func drawVPatternOnImage(img *ebiten.Image, size int, clr color.RGBA) {
	cx := float64(size) / 2
	thickness := 3.0
	for py := 0; py < size; py++ {
		t := float64(py) / float64(size)
		leftX := cx - cx*0.6*(1-t)
		rightX := cx + cx*0.6*(1-t)
		for px := 0; px < size; px++ {
			fPx := float64(px)
			dLeft := fPx - leftX
			dRight := fPx - rightX
			if dLeft*dLeft < thickness*thickness || dRight*dRight < thickness*thickness {
				img.Set(px, py, clr)
			}
		}
	}
}

// --- Public drawing utilities ---

// DrawSprite draws a pre-built sprite image onto dst at (x,y) top-left.
func DrawSprite(dst *ebiten.Image, sprite *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(sprite, op)
}

// DrawSpriteCentered draws a sprite centered at (cx, cy) in pixel coordinates.
func DrawSpriteCentered(dst *ebiten.Image, sprite *ebiten.Image, cx, cy float64) {
	w := float64(sprite.Bounds().Dx())
	h := float64(sprite.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx-w/2, cy-h/2)
	dst.DrawImage(sprite, op)
}

// DrawFilledRect draws a filled rectangle onto dst at (x,y) with given size and color.
func DrawFilledRect(dst *ebiten.Image, x, y, w, h float64, clr color.RGBA) {
	iw := int(w)
	ih := int(h)
	if iw <= 0 || ih <= 0 {
		return
	}
	key := fmt.Sprintf("rect_%d_%d_%d_%d_%d_%d", iw, ih, clr.R, clr.G, clr.B, clr.A)
	img := getCachedSprite(key, iw, ih, func(img *ebiten.Image) {
		img.Fill(clr)
	})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(img, op)
}

// DrawFilledCircle draws a filled circle onto dst centered at (cx,cy).
func DrawFilledCircle(dst *ebiten.Image, cx, cy, radius float64, clr color.RGBA) {
	iRadius := int(radius)
	if iRadius <= 0 {
		return
	}
	size := iRadius*2 + 2
	key := fmt.Sprintf("circle_%d_%d_%d_%d_%d", iRadius, clr.R, clr.G, clr.B, clr.A)
	img := getCachedSprite(key, size, size, func(img *ebiten.Image) {
		drawCircleOnImage(img, radius, clr)
	})
	center := float64(size) / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx-center, cy-center)
	dst.DrawImage(img, op)
}

// DrawCross draws a plus/cross shape onto dst centered at (cx,cy).
func DrawCross(dst *ebiten.Image, cx, cy, size float64, clr color.RGBA) {
	iSize := int(size)
	if iSize <= 0 {
		return
	}
	key := fmt.Sprintf("cross_%d_%d_%d_%d_%d", iSize, clr.R, clr.G, clr.B, clr.A)
	img := getCachedSprite(key, iSize, iSize, func(img *ebiten.Image) {
		drawCrossOnImage(img, iSize, clr)
	})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx-float64(iSize)/2, cy-float64(iSize)/2)
	dst.DrawImage(img, op)
}

// ClearSpriteCache removes all cached sprites.
func ClearSpriteCache() {
	spriteCacheMu.Lock()
	defer spriteCacheMu.Unlock()
	spriteCache = make(map[string]*ebiten.Image)
}
