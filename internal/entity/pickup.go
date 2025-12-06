package entity

import "github.com/arkosh/vampyre-survival/internal/physics"

type PickupType int

const (
	PickupXP PickupType = iota
	PickupHP
)

const MagnetSpeed = 12.0

type Pickup struct {
	Body      physics.Body
	Type      PickupType
	Value     int
	Collected bool
}

func NewXPGem(x, y float64, value int) *Pickup {
	return &Pickup{
		Body: physics.Body{
			Pos:   physics.Vec2{X: x, Y: y},
			Width: 1, Height: 1,
		},
		Type:  PickupXP,
		Value: value,
	}
}

func (p *Pickup) MagnetToward(target physics.Vec2, magnetRadius float64) {
	dist := p.Body.Pos.DistanceTo(target)
	if dist > magnetRadius || dist < 0.1 {
		return
	}
	dir := target.Sub(p.Body.Pos).Normalize()
	p.Body.Vel = dir.Scale(MagnetSpeed)
}

func (p *Pickup) Update(dt float64) {
	p.Body.Update(dt)
}

func (p *Pickup) Collect() {
	p.Collected = true
}

func NewHPPickup(x, y float64, value int) *Pickup {
	return &Pickup{
		Body: physics.Body{
			Pos:   physics.Vec2{X: x, Y: y},
			Width: 1, Height: 1,
		},
		Type:  PickupHP,
		Value: value,
	}
}
