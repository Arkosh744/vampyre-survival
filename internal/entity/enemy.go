package entity

import "github.com/arkosh/vampyre-survival/internal/physics"

type EnemyType int

const (
	EnemyNormal EnemyType = iota
	EnemyBoss
	EnemySwarmer
	EnemyTank
	EnemyDasher
)

const (
	DasherIdleSpeed  = 3.0
	DasherDashSpeed  = 20.0
	DasherDashRange  = 15.0
)

type Enemy struct {
	Body    physics.Body
	Type    EnemyType
	HP      int
	MaxHP   int
	Damage  int
	Speed   float64
	XPDrop  int
	Dashing bool // for Dasher telegraph
}

func NewEnemy(typ EnemyType, x, y float64) *Enemy {
	e := &Enemy{Type: typ}

	switch typ {
	case EnemyBoss:
		e.HP = 80
		e.MaxHP = 80
		e.Damage = 20
		e.Speed = 6.0
		e.XPDrop = 10
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 6.0,
			Width:    3,
			Height:   3,
		}
	case EnemySwarmer:
		e.HP = 5
		e.MaxHP = 5
		e.Damage = 3
		e.Speed = 8.0
		e.XPDrop = 1
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 8.0,
			Width:    1,
			Height:   1,
		}
	case EnemyTank:
		e.HP = 60
		e.MaxHP = 60
		e.Damage = 15
		e.Speed = 2.0
		e.XPDrop = 5
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 2.0,
			Width:    2,
			Height:   2,
		}
	case EnemyDasher:
		e.HP = 10
		e.MaxHP = 10
		e.Damage = 25
		e.Speed = DasherIdleSpeed
		e.XPDrop = 3
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: DasherDashSpeed,
			Width:    1,
			Height:   1,
		}
	default: // EnemyNormal
		e.HP = 15
		e.MaxHP = 15
		e.Damage = 10
		e.Speed = 4.5
		e.XPDrop = 1
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 4.5,
			Width:    1,
			Height:   1,
		}
	}

	return e
}

func (e *Enemy) ChaseTarget(target physics.Vec2) {
	dir := target.Sub(e.Body.Pos).Normalize()

	if e.Type == EnemyDasher {
		dist := e.Body.Pos.DistanceTo(target)
		if dist < DasherDashRange {
			e.Speed = DasherDashSpeed
			e.Dashing = true
		} else {
			e.Speed = DasherIdleSpeed
			e.Dashing = false
		}
	}

	e.Body.Vel = dir.Scale(e.Speed)
}

func (e *Enemy) Update(dt float64) {
	e.Body.Update(dt)
}

func (e *Enemy) TakeDamage(dmg int) bool {
	e.HP -= dmg
	if e.HP <= 0 {
		e.HP = 0
		return true
	}
	return false
}

func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}

// ScaleForWave increases HP (+8%/wave) and Damage (+5%/wave). Speed unchanged.
func (e *Enemy) ScaleForWave(wave int) {
	if wave <= 1 {
		return
	}
	hpMult := 1.0 + 0.08*float64(wave-1)
	dmgMult := 1.0 + 0.05*float64(wave-1)
	e.HP = int(float64(e.HP) * hpMult)
	e.MaxHP = e.HP
	e.Damage = int(float64(e.Damage) * dmgMult)
}
