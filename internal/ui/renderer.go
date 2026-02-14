package ui

import (
	"fmt"
	"math"
	"strconv"

	"github.com/arkosh/vampyre-survival/internal/engine"
	"github.com/arkosh/vampyre-survival/internal/entity"
	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/arkosh/vampyre-survival/internal/weapon"
	"github.com/hajimehoshi/ebiten/v2"
)

type VisualEffect struct {
	Pos   physics.Vec2
	Char  rune
	Color string
	TTL   float64
}

type TextEffect struct {
	Pos   physics.Vec2
	Text  string
	Color string
	TTL   float64
	VelY  float64
}

type Renderer struct {
	Camera      *engine.Camera
	Effects     []VisualEffect
	TextEffects []TextEffect
}

func NewRenderer(cam *engine.Camera) *Renderer {
	return &Renderer{Camera: cam}
}

func (r *Renderer) AddDamageNumber(pos physics.Vec2, damage int) {
	r.AddDamageNumberColored(pos, damage, ColorYellow)
}

func (r *Renderer) AddDamageNumberColored(pos physics.Vec2, damage int, color string) {
	r.TextEffects = append(r.TextEffects, TextEffect{
		Pos:   pos,
		Text:  strconv.Itoa(damage),
		Color: color,
		TTL:   0.5,
		VelY:  -3,
	})
}

func (r *Renderer) UpdateEffects(dt float64) {
	alive := r.Effects[:0]
	for _, e := range r.Effects {
		e.TTL -= dt
		if e.TTL > 0 {
			alive = append(alive, e)
		}
	}
	r.Effects = alive

	aliveText := r.TextEffects[:0]
	for _, te := range r.TextEffects {
		te.TTL -= dt
		te.Pos.Y += te.VelY * dt
		if te.TTL > 0 {
			aliveText = append(aliveText, te)
		}
	}
	r.TextEffects = aliveText
}

// DrawHUD renders the heads-up display directly onto an *ebiten.Image
// using pixel coordinates.
func (r *Renderer) DrawHUD(dst *ebiten.Image, p *entity.Player, wave, totalWaves, kills int, weapons []weapon.WeaponStats, comboCount int, comboDmgPct float64) {
	screenW := dst.Bounds().Dx()

	// HP bar: top-left
	hpBarX := 5.0
	hpBarY := 4.0
	hpBarW := 100.0
	hpBarH := 8.0

	// Background (dark)
	DrawFilledRect(dst, hpBarX, hpBarY, hpBarW, hpBarH, ColorToRGBA(ColorDim))
	// Filled portion
	hpFrac := float64(p.HP) / float64(p.MaxHP)
	if hpFrac < 0 {
		hpFrac = 0
	}
	if hpFrac > 1 {
		hpFrac = 1
	}
	DrawFilledRect(dst, hpBarX, hpBarY, hpBarW*hpFrac, hpBarH, ColorToRGBA(ColorRed))
	// Label
	hpLabel := fmt.Sprintf("HP: %d/%d", p.HP, p.MaxHP)
	DrawText(dst, int(hpBarX)+2, int(hpBarY)+0, hpLabel, ColorToRGBA(ColorWhite))

	// XP bar: below HP
	xpBarY := hpBarY + hpBarH + 2
	xpBarH := 6.0
	DrawFilledRect(dst, hpBarX, xpBarY, hpBarW, xpBarH, ColorToRGBA(ColorDim))
	xpFrac := float64(p.XP) / float64(p.XPToNextLevel())
	if xpFrac > 1 {
		xpFrac = 1
	}
	DrawFilledRect(dst, hpBarX, xpBarY, hpBarW*xpFrac, xpBarH, ColorToRGBA(ColorCyan))
	xpLabel := fmt.Sprintf("XP Lv.%d", p.Level)
	DrawText(dst, int(hpBarX)+2, int(xpBarY)+1, xpLabel, ColorToRGBA(ColorWhite))

	// Weapon info: below XP bar, up to 3 weapon slots
	weaponColors := map[weapon.WeaponKind]string{
		weapon.KindSword:      ColorWhite,
		weapon.KindProjectile: ColorCyan,
		weapon.KindAoE:        ColorMagenta,
		weapon.KindLightning:  ColorCyan,
		weapon.KindOrbital:    ColorMagenta,
	}
	for i, ws := range weapons {
		if i >= 3 {
			break
		}
		wy := xpBarY + xpBarH + 4 + float64(i)*14.0
		line := fmt.Sprintf("%-6s Lv.%d DMG:%d", ws.Name, ws.Level, ws.Damage)
		wColor := weaponColors[ws.Kind]
		if wColor == "" {
			wColor = ColorWhite
		}
		DrawText(dst, int(hpBarX), int(wy), line, ColorToRGBA(wColor))
	}

	// Combo indicator: top-right corner when combo >= 5
	if comboCount >= 5 {
		comboStr := fmt.Sprintf("COMBO: x%d (+%.0f%%)", comboCount, comboDmgPct)
		comboColor := ColorYellow
		if comboCount >= 20 {
			comboColor = ColorBoldRed
		} else if comboCount >= 10 {
			comboColor = ColorMagenta
		}
		// Approximate 7px per character width for basic font
		textW := len(comboStr) * 7
		DrawText(dst, screenW-textW-5, 4, comboStr, ColorToRGBA(comboColor))
	}
}

// DrawBottomBar renders wave, enemies remaining, and kills at the bottom.
func (r *Renderer) DrawBottomBar(dst *ebiten.Image, wave, totalWaves, kills, enemiesLeft int) {
	screenW := dst.Bounds().Dx()
	screenH := dst.Bounds().Dy()
	barH := 14.0
	barY := float64(screenH) - barH

	DrawFilledRect(dst, 0, barY, float64(screenW), barH, ColorToRGBA(ColorDim))

	waveStr := fmt.Sprintf("Wave %d/%d", wave, totalWaves)
	DrawText(dst, 5, int(barY)+1, waveStr, ColorToRGBA(ColorWhite))

	enemyStr := fmt.Sprintf("Enemies: %d", enemiesLeft)
	enemyW := len(enemyStr) * 7
	DrawText(dst, (screenW-enemyW)/2, int(barY)+1, enemyStr, ColorToRGBA(ColorYellow))

	killStr := fmt.Sprintf("Kills: %d", kills)
	killW := len(killStr) * 7
	DrawText(dst, screenW-killW-5, int(barY)+1, killStr, ColorToRGBA(ColorWhite))
}

// DrawGround fills the screen with grass tiles that scroll with the camera.
func (r *Renderer) DrawGround(dst *ebiten.Image) {
	screenW := dst.Bounds().Dx()
	screenH := dst.Bounds().Dy()
	ts := float64(engine.TileSize)

	// Calculate visible tile range
	startX := int(math.Floor(r.Camera.Pos.X)) - 1
	startY := int(math.Floor(r.Camera.Pos.Y)) - 1
	endX := int(math.Ceil(r.Camera.Pos.X+float64(screenW)/ts)) + 1
	endY := int(math.Ceil(r.Camera.Pos.Y+float64(screenH)/ts)) + 1

	for wy := startY; wy <= endY; wy++ {
		for wx := startX; wx <= endX; wx++ {
			variant := ((wx*7 + wy*13) % 4)
			if variant < 0 {
				variant += 4
			}
			sprite := GrassSprite(variant)
			px := (float64(wx) - r.Camera.Pos.X) * ts
			py := (float64(wy) - r.Camera.Pos.Y) * ts
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(px, py)
			dst.DrawImage(sprite, op)
		}
	}
}

// DrawWorldBorder draws the world border as 2px wide colored lines along
// the world edges, only rendering visible portions.
func (r *Renderer) DrawWorldBorder(dst *ebiten.Image, worldW, worldH float64) {
	borderClr := ColorToRGBA(ColorWhite)
	ts := float64(engine.TileSize)

	// Top border
	for x := 0.0; x < worldW; x++ {
		pos := physics.Vec2{X: x, Y: 0}
		if r.Camera.IsVisible(pos) {
			px, py := r.Camera.WorldToScreenPx(pos)
			DrawFilledRect(dst, px, py, ts, 2, borderClr)
		}
	}
	// Bottom border
	for x := 0.0; x < worldW; x++ {
		pos := physics.Vec2{X: x, Y: worldH - 1}
		if r.Camera.IsVisible(pos) {
			px, py := r.Camera.WorldToScreenPx(pos)
			DrawFilledRect(dst, px, py+ts-2, ts, 2, borderClr)
		}
	}
	// Left border
	for y := 0.0; y < worldH; y++ {
		pos := physics.Vec2{X: 0, Y: y}
		if r.Camera.IsVisible(pos) {
			px, py := r.Camera.WorldToScreenPx(pos)
			DrawFilledRect(dst, px, py, 2, ts, borderClr)
		}
	}
	// Right border
	for y := 0.0; y < worldH; y++ {
		pos := physics.Vec2{X: worldW - 1, Y: y}
		if r.Camera.IsVisible(pos) {
			px, py := r.Camera.WorldToScreenPx(pos)
			DrawFilledRect(dst, px+ts-2, py, 2, ts, borderClr)
		}
	}
}

// DrawPlayer renders the player as a pixel sprite centered at their
// world position, with a direction indicator.
func (r *Renderer) DrawPlayer(dst *ebiten.Image, p *entity.Player) {
	if !r.Camera.IsVisible(p.Body.Pos) {
		return
	}
	px, py := r.Camera.WorldToScreenPx(p.Body.Pos)
	sprite := PlayerSprite(p.Invulnerable)
	DrawSpriteCentered(dst, sprite, px, py)

	// Direction indicator
	if p.Direction.Length() > 0 {
		dir := p.Direction.Normalize()
		indicator := PlayerDirectionIndicator(p.Invulnerable)
		dx := dir.X * float64(engine.TileSize)
		dy := dir.Y * float64(engine.TileSize)
		DrawSpriteCentered(dst, indicator, px+dx, py+dy)
	}
}

// DrawEnemy renders an enemy as a pixel sprite appropriate for its type.
func (r *Renderer) DrawEnemy(dst *ebiten.Image, e *entity.Enemy) {
	if !r.Camera.IsVisible(e.Body.Pos) {
		return
	}
	px, py := r.Camera.WorldToScreenPx(e.Body.Pos)

	var spriteType string
	var variant bool
	switch e.Type {
	case entity.EnemyElderVampyre:
		if e.Enraged {
			spriteType = SpriteEnemyElderEnraged
		} else {
			spriteType = SpriteEnemyElder
		}
	case entity.EnemyBoss:
		spriteType = SpriteEnemyBoss
	case entity.EnemySwarmer:
		spriteType = SpriteEnemySwarmer
	case entity.EnemyTank:
		spriteType = SpriteEnemyTank
	case entity.EnemyDasher:
		spriteType = SpriteEnemyDasher
		variant = e.Dashing
	default:
		spriteType = SpriteEnemyNormal
	}

	sprite := EnemySprite(spriteType, variant)
	DrawSpriteCentered(dst, sprite, px, py)
}

// DrawPickup renders a pickup item as a pixel sprite.
func (r *Renderer) DrawPickup(dst *ebiten.Image, p *entity.Pickup) {
	if p.Collected || !r.Camera.IsVisible(p.Body.Pos) {
		return
	}
	px, py := r.Camera.WorldToScreenPx(p.Body.Pos)

	var pickupType string
	switch p.Type {
	case entity.PickupHP:
		pickupType = SpritePickupHP
	default:
		pickupType = SpritePickupXP
	}
	sprite := PickupSprite(pickupType, p.Value)
	DrawSpriteCentered(dst, sprite, px, py)
}

// DrawWeaponVisuals renders weapon visuals as pixel sprites. Visuals with
// TTL > 0 are added to the effects list; others are drawn immediately.
func (r *Renderer) DrawWeaponVisuals(dst *ebiten.Image, visuals []weapon.Visual) {
	for _, v := range visuals {
		if v.TTL > 0 {
			r.Effects = append(r.Effects, VisualEffect{
				Pos: v.Pos, Char: v.Char, Color: ColorYellow, TTL: v.TTL,
			})
			continue
		}
		if !r.Camera.IsVisible(v.Pos) {
			continue
		}
		px, py := r.Camera.WorldToScreenPx(v.Pos)
		var kind string
		switch v.Char {
		case 'z', '~':
			kind = "lightning"
		case '\u25cf': // bullet
			kind = "orbital"
		default:
			kind = "default"
		}
		sprite := WeaponEffectSprite(kind)
		DrawSpriteCentered(dst, sprite, px, py)
	}
}

// DrawEffects renders death particles and floating text effects as pixel
// graphics.
func (r *Renderer) DrawEffects(dst *ebiten.Image) {
	for _, e := range r.Effects {
		if !r.Camera.IsVisible(e.Pos) {
			continue
		}
		px, py := r.Camera.WorldToScreenPx(e.Pos)
		// Death particles as pixel-art star bursts
		sprite := DeathParticleSprite(ColorToRGBA(e.Color))
		DrawSpriteCentered(dst, sprite, px, py)
	}
	for _, te := range r.TextEffects {
		if !r.Camera.IsVisible(te.Pos) {
			continue
		}
		px, py := r.Camera.WorldToScreenPx(te.Pos)
		DrawText(dst, int(px), int(py), te.Text, ColorToRGBA(te.Color))
	}
}

// DrawBossHP renders a boss health bar centered at the top of the screen.
func (r *Renderer) DrawBossHP(dst *ebiten.Image, name string, hp, maxHP int) {
	screenW := dst.Bounds().Dx()
	barW := 150.0
	barH := 10.0
	barX := (float64(screenW) - barW) / 2
	barY := 4.0

	// Background track
	DrawFilledRect(dst, barX, barY, barW, barH, ColorToRGBA(ColorDim))
	// Filled portion
	frac := float64(hp) / float64(maxHP)
	if frac < 0 {
		frac = 0
	}
	DrawFilledRect(dst, barX, barY, barW*frac, barH, ColorToRGBA(ColorBoldRed))
	// Centered label
	label := fmt.Sprintf("%s %d/%d", name, hp, maxHP)
	textW := len(label) * 7
	DrawText(dst, (screenW-textW)/2, int(barY)+0, label, ColorToRGBA(ColorWhite))
}
