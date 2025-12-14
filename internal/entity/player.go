package entity

import (
	"math"

	"github.com/arkosh/vampyre-survival/internal/physics"
)

const (
	PlayerSpeed       = 12.0
	PlayerMaxHP       = 100
	PlayerWidth       = 1.0
	PlayerHeight      = 1.0
	InvulnDuration    = 0.5
	BaseXPToLevel     = 5.0
	XPLevelMultiplier = 1.3
)

type Player struct {
	Body         physics.Body
	HP           int
	MaxHP        int
	XP           int
	Level        int
	Kills        int
	Direction    physics.Vec2
	Invulnerable bool
	InvulnTimer  float64
	InvulnBonus  float64
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		Body: physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			Friction: 0.85,
			MaxSpeed: PlayerSpeed,
			Width:    PlayerWidth,
			Height:   PlayerHeight,
		},
		HP:    PlayerMaxHP,
		MaxHP: PlayerMaxHP,
		Level: 1,
	}
}

func (p *Player) SetDirection(dir physics.Vec2) {
	p.Direction = dir
	if dir.Length() > 0 {
		p.Body.ApplyForce(dir.Normalize().Scale(PlayerSpeed * 0.5))
	}
}

func (p *Player) Update(dt float64) {
	p.Body.Update(dt)

	if p.Invulnerable {
		p.InvulnTimer -= dt
		if p.InvulnTimer <= 0 {
			p.Invulnerable = false
			p.InvulnTimer = 0
		}
	}
}

func (p *Player) TakeDamage(dmg int) bool {
	if p.Invulnerable {
		return false
	}
	p.HP -= dmg
	if p.HP <= 0 {
		p.HP = 0
		return true
	}
	p.Invulnerable = true
	p.InvulnTimer = InvulnDuration + p.InvulnBonus
	return false
}

func (p *Player) AddXP(amount int) bool {
	p.XP += amount
	needed := p.XPToNextLevel()
	if p.XP >= needed {
		p.XP -= needed
		p.Level++
		return true
	}
	return false
}

func (p *Player) XPToNextLevel() int {
	return int(BaseXPToLevel * math.Pow(XPLevelMultiplier, float64(p.Level-1)))
}

func (p *Player) AddMaxHP(v int) {
	p.MaxHP += v
	p.HP += v
}

func (p *Player) AddSpeed(v float64) {
	p.Body.MaxSpeed += v
}
