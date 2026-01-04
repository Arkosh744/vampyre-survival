package ui

import (
	"fmt"
	"strconv"

	"github.com/arkosh/vampyre-survival/internal/engine"
	"github.com/arkosh/vampyre-survival/internal/entity"
	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/arkosh/vampyre-survival/internal/weapon"
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
	Term        *Terminal
	Camera      *engine.Camera
	Effects     []VisualEffect
	TextEffects []TextEffect
}

func NewRenderer(t *Terminal, cam *engine.Camera) *Renderer {
	return &Renderer{Term: t, Camera: cam}
}

// DrawGround draws a sparse dot pattern that scrolls with the camera,
// giving the player a strong visual reference for movement.
func (r *Renderer) DrawGround() {
	camX := int(r.Camera.Pos.X)
	camY := int(r.Camera.Pos.Y)
	viewH := r.Camera.Height
	if viewH > r.Term.Height()-6 {
		viewH = r.Term.Height() - 6
	}

	for sy := 0; sy < viewH; sy++ {
		for sx := 0; sx < r.Camera.Width; sx++ {
			wx := camX + sx
			wy := camY + sy
			// Deterministic pseudo-random pattern: dots every ~5 cells
			if (wx*7+wy*13)%23 == 0 {
				r.Term.SetCell(sx, sy, '.', ColorDim)
			}
		}
	}
}

func (r *Renderer) DrawPlayer(p *entity.Player) {
	if !r.Camera.IsVisible(p.Body.Pos) {
		return
	}
	sx, sy := r.Camera.WorldToScreen(p.Body.Pos)
	color := ColorGreen
	if p.Invulnerable {
		color = ColorYellow
	}
	r.Term.SetCell(sx, sy, '@', color)

	// Draw direction indicator around player
	if p.Direction.Length() > 0 {
		dir := p.Direction.Normalize()
		dx := int(dir.X + 0.5)
		dy := int(dir.Y + 0.5)
		if dx == 0 && dy == 0 {
			if dir.X > 0 {
				dx = 1
			} else if dir.X < 0 {
				dx = -1
			}
		}
		r.Term.SetCell(sx+dx, sy+dy, '+', color)
	}
}

func (r *Renderer) DrawEnemy(e *entity.Enemy) {
	if !r.Camera.IsVisible(e.Body.Pos) {
		return
	}
	sx, sy := r.Camera.WorldToScreen(e.Body.Pos)
	switch e.Type {
	case entity.EnemyElderVampyre:
		color := ColorMagenta
		if e.Enraged {
			color = ColorBoldRed
		}
		for dy := -1; dy <= 1; dy++ {
			for dx := -2; dx <= 2; dx++ {
				ch := rune('E')
				if dx == 0 && dy == 0 {
					ch = 'V'
				}
				r.Term.SetCell(sx+dx, sy+dy, ch, color)
			}
		}
	case entity.EnemyBoss:
		r.Term.SetCell(sx, sy, 'B', ColorBoldRed)
	case entity.EnemySwarmer:
		r.Term.SetCell(sx, sy, 's', ColorYellow)
	case entity.EnemyTank:
		r.Term.SetCell(sx, sy, 'T', ColorMagenta)
	case entity.EnemyDasher:
		ch := rune('d')
		if e.Dashing {
			ch = 'D'
		}
		r.Term.SetCell(sx, sy, ch, ColorCyan)
	default:
		r.Term.SetCell(sx, sy, 'V', ColorRed)
	}
}

func (r *Renderer) DrawPickup(p *entity.Pickup) {
	if p.Collected || !r.Camera.IsVisible(p.Body.Pos) {
		return
	}
	sx, sy := r.Camera.WorldToScreen(p.Body.Pos)
	switch p.Type {
	case entity.PickupHP:
		r.Term.SetCell(sx, sy, '+', ColorGreen)
	default:
		if p.Value >= 10 {
			r.Term.SetCell(sx, sy, '*', ColorCyan)
		} else {
			r.Term.SetCell(sx, sy, 'o', ColorBlue)
		}
	}
}

func (r *Renderer) DrawWeaponVisuals(visuals []weapon.Visual) {
	for _, v := range visuals {
		color := ColorYellow
		switch v.Char {
		case 'z', '~':
			color = ColorCyan
		case '●':
			color = ColorMagenta
		}
		if v.TTL <= 0 {
			if r.Camera.IsVisible(v.Pos) {
				sx, sy := r.Camera.WorldToScreen(v.Pos)
				r.Term.SetCell(sx, sy, v.Char, color)
			}
		} else {
			r.Effects = append(r.Effects, VisualEffect{
				Pos: v.Pos, Char: v.Char, Color: color, TTL: v.TTL,
			})
		}
	}
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

func (r *Renderer) DrawEffects() {
	for _, e := range r.Effects {
		if !r.Camera.IsVisible(e.Pos) {
			continue
		}
		sx, sy := r.Camera.WorldToScreen(e.Pos)
		r.Term.SetCell(sx, sy, e.Char, e.Color)
	}
	for _, te := range r.TextEffects {
		if !r.Camera.IsVisible(te.Pos) {
			continue
		}
		sx, sy := r.Camera.WorldToScreen(te.Pos)
		r.Term.WriteStr(sx, sy, te.Text, te.Color)
	}
}

func (r *Renderer) DrawWorldBorder(worldW, worldH float64) {
	for x := 0.0; x < worldW; x++ {
		r.drawBorderCell(physics.Vec2{X: x, Y: 0})
		r.drawBorderCell(physics.Vec2{X: x, Y: worldH - 1})
	}
	for y := 0.0; y < worldH; y++ {
		r.drawBorderCell(physics.Vec2{X: 0, Y: y})
		r.drawBorderCell(physics.Vec2{X: worldW - 1, Y: y})
	}
}

func (r *Renderer) drawBorderCell(pos physics.Vec2) {
	if !r.Camera.IsVisible(pos) {
		return
	}
	sx, sy := r.Camera.WorldToScreen(pos)
	r.Term.SetCell(sx, sy, '#', ColorWhite)
}

func (r *Renderer) DrawHUD(p *entity.Player, wave int, totalWaves int, kills int, weapons []weapon.WeaponStats, comboCount int, comboDmgPct float64) {
	h := r.Term.Height()
	w := r.Term.Width()

	hpBar := renderBar(p.HP, p.MaxHP, 20)
	xpBar := renderBar(p.XP, p.XPToNextLevel(), 20)

	// Row h-6: HP + combo
	hpLine := fmt.Sprintf("HP: %s %d/%d", hpBar, p.HP, p.MaxHP)
	r.Term.WriteStr(0, h-6, hpLine, ColorRed)

	if comboCount >= 5 {
		comboStr := fmt.Sprintf("COMBO: x%d (+%.2f%%)", comboCount, comboDmgPct)
		color := ColorYellow
		if comboCount >= 20 {
			color = ColorBoldRed
		} else if comboCount >= 10 {
			color = ColorMagenta
		}
		r.Term.WriteStr(w-len(comboStr), h-6, comboStr, color)
	}

	// Row h-5: XP
	r.Term.WriteStr(0, h-5, fmt.Sprintf("XP: %s Lv.%d", xpBar, p.Level), ColorCyan)

	// Rows h-4..h-2: weapon stats (up to 3 slots)
	weaponColors := map[weapon.WeaponKind]string{
		weapon.KindSword:      ColorWhite,
		weapon.KindProjectile: ColorCyan,
		weapon.KindAoE:        ColorMagenta,
		weapon.KindLightning:  ColorCyan,
		weapon.KindOrbital:    ColorMagenta,
	}
	for i := 0; i < 3; i++ {
		row := h - 4 + i
		if i < len(weapons) {
			ws := weapons[i]
			line := fmt.Sprintf("%-6s Lv.%d  DMG:%d", ws.Name, ws.Level, ws.Damage)
			if ws.Cooldown > 0 {
				line += fmt.Sprintf("  CD:%.2f", ws.Cooldown)
			}
			if ws.Range > 0 {
				line += fmt.Sprintf("  RNG:%.1f", ws.Range)
			}
			color := weaponColors[ws.Kind]
			if color == "" {
				color = ColorWhite
			}
			r.Term.WriteStr(0, row, line, color)
		}
	}

	// Row h-1: wave info
	r.Term.WriteStr(0, h-1, fmt.Sprintf("Wave: %d/%d  Kills: %d", wave, totalWaves, kills), ColorWhite)
}

func (r *Renderer) DrawBossHP(name string, hp, maxHP int) {
	w := r.Term.Width()
	bar := renderBar(hp, maxHP, 30)
	line := fmt.Sprintf("%s %s %d/%d", name, bar, hp, maxHP)
	r.Term.WriteStr((w-len(line))/2, 1, line, ColorBoldRed)
}

func renderBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	filled := current * width / max
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	bar := make([]byte, width)
	for i := 0; i < width; i++ {
		if i < filled {
			bar[i] = '='
		} else {
			bar[i] = '-'
		}
	}
	return "[" + string(bar) + "]"
}
