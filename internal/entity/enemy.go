package entity

import (
	"math"

	"github.com/arkosh/vampyre-survival/internal/physics"
)

type EnemyType int

const (
	EnemyNormal EnemyType = iota
	EnemyBoss
	EnemySwarmer
	EnemyTank
	EnemyDasher
	EnemyElderVampyre
)

const (
	DasherIdleSpeed  = 6.0
	DasherDashSpeed  = 40.0
	DasherDashRange  = 15.0

	ElderBaseHP      = 800
	ElderBaseDMG     = 30
	ElderBaseSpeed   = 5.0
	ElderEnrageHP    = 0.5
	ElderEnrageSpeed = 8.0
	ElderEnrageDMG   = 45
	ElderSpawnCD     = 3.0
	ElderEnrageCD    = 1.5
	ElderSpawnCount  = 3
	ElderEnrageCount = 5
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

	// Elder Vampyre fields
	SpawnTimer float64
	SpawnCD    float64
	Enraged    bool
}

func NewEnemy(typ EnemyType, x, y float64) *Enemy {
	e := &Enemy{Type: typ}

	switch typ {
	case EnemyBoss:
		e.HP = 80
		e.MaxHP = 80
		e.Damage = 20
		e.Speed = 12.0
		e.XPDrop = 5
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 12.0,
			Width:    3,
			Height:   3,
		}
	case EnemySwarmer:
		e.HP = 5
		e.MaxHP = 5
		e.Damage = 3
		e.Speed = 16.0
		e.XPDrop = 1
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 16.0,
			Width:    1,
			Height:   1,
		}
	case EnemyTank:
		e.HP = 60
		e.MaxHP = 60
		e.Damage = 15
		e.Speed = 4.0
		e.XPDrop = 3
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 4.0,
			Width:    2,
			Height:   2,
		}
	case EnemyDasher:
		e.HP = 10
		e.MaxHP = 10
		e.Damage = 25
		e.Speed = DasherIdleSpeed
		e.XPDrop = 2
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: DasherDashSpeed,
			Width:    1,
			Height:   1,
		}
	case EnemyElderVampyre:
		e.HP = ElderBaseHP
		e.MaxHP = ElderBaseHP
		e.Damage = ElderBaseDMG
		e.Speed = ElderBaseSpeed
		e.XPDrop = 25
		e.SpawnCD = ElderSpawnCD
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: ElderBaseSpeed,
			Width:    5,
			Height:   3,
		}
	default: // EnemyNormal
		e.HP = 15
		e.MaxHP = 15
		e.Damage = 10
		e.Speed = 9.0
		e.XPDrop = 1
		e.Body = physics.Body{
			Pos:      physics.Vec2{X: x, Y: y},
			MaxSpeed: 9.0,
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

// CheckEnrage triggers phase 2 at ≤50% HP. Returns true the first time it triggers.
func (e *Enemy) CheckEnrage() bool {
	if e.Type != EnemyElderVampyre || e.Enraged {
		return false
	}
	if float64(e.HP) <= float64(e.MaxHP)*ElderEnrageHP {
		e.Enraged = true
		e.Speed = ElderEnrageSpeed
		e.Body.MaxSpeed = ElderEnrageSpeed
		e.Damage = ElderEnrageDMG
		e.SpawnCD = ElderEnrageCD
		return true
	}
	return false
}

// UpdateSpawnTimer ticks the spawn timer. Returns true when cooldown reached.
func (e *Enemy) UpdateSpawnTimer(dt float64) bool {
	if e.SpawnCD <= 0 {
		return false
	}
	e.SpawnTimer += dt
	if e.SpawnTimer >= e.SpawnCD {
		e.SpawnTimer -= e.SpawnCD
		return true
	}
	return false
}

// SpawnCount returns how many minions to spawn per cycle.
func (e *Enemy) SpawnCount() int {
	if e.Enraged {
		return ElderEnrageCount
	}
	return ElderSpawnCount
}

// ScaleForWave increases HP (x1.12/wave), Damage (x1.08/wave), Speed (x1.075/wave).
func (e *Enemy) ScaleForWave(wave int) {
	if wave <= 1 {
		return
	}
	hpMult := math.Pow(1.12, float64(wave-1))
	dmgMult := math.Pow(1.08, float64(wave-1))
	spdMult := math.Pow(1.075, float64(wave-1))
	e.HP = int(float64(e.HP) * hpMult)
	e.MaxHP = e.HP
	e.Damage = int(float64(e.Damage) * dmgMult)
	e.Speed *= spdMult
	e.Body.MaxSpeed *= spdMult
}
