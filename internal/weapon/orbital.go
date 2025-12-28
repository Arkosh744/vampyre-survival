package weapon

import (
	"math"

	"github.com/arkosh/vampyre-survival/internal/physics"
)

const (
	OrbitalBaseDamage    = 12
	OrbitalBaseRadius    = 4.0
	OrbitalBaseOrbCount  = 2
	OrbitalRotationSpeed = 3.0
	OrbitalHitRadius     = 1.5
	OrbitalHitCooldown   = 0.5
)

type Orbital struct {
	level         int
	damage        int
	radius        float64
	orbCount      int
	angle         float64
	rotationSpeed float64
	hitCooldowns  map[int]float64
	hits          []HitResult
	visuals       []Visual
}

func NewOrbital() *Orbital {
	return &Orbital{
		level:         1,
		damage:        OrbitalBaseDamage,
		radius:        OrbitalBaseRadius,
		orbCount:      OrbitalBaseOrbCount,
		rotationSpeed: OrbitalRotationSpeed,
		hitCooldowns:  make(map[int]float64),
	}
}

func (o *Orbital) Update(dt float64, ownerPos physics.Vec2, targets []Target) {
	o.hits = nil
	o.visuals = nil

	o.angle += o.rotationSpeed * dt

	// Decay hit cooldowns
	for id, cd := range o.hitCooldowns {
		cd -= dt
		if cd <= 0 {
			delete(o.hitCooldowns, id)
		} else {
			o.hitCooldowns[id] = cd
		}
	}

	for i := 0; i < o.orbCount; i++ {
		a := o.angle + 2*math.Pi*float64(i)/float64(o.orbCount)
		orbPos := physics.Vec2{
			X: ownerPos.X + math.Cos(a)*o.radius,
			Y: ownerPos.Y + math.Sin(a)*o.radius,
		}

		o.visuals = append(o.visuals, Visual{Pos: orbPos, Char: '●', TTL: 0})

		for _, t := range targets {
			if orbPos.DistanceTo(t.Pos) <= OrbitalHitRadius {
				if _, onCD := o.hitCooldowns[t.ID]; !onCD {
					o.hits = append(o.hits, HitResult{TargetID: t.ID, Damage: o.damage})
					o.hitCooldowns[t.ID] = OrbitalHitCooldown
				}
			}
		}
	}
}

func (o *Orbital) GetHits() []HitResult { return o.hits }
func (o *Orbital) GetVisuals() []Visual  { return o.visuals }
func (o *Orbital) Level() int            { return o.level }

func (o *Orbital) Upgrade() {
	o.level++
	o.damage += 3
	o.radius += 0.5
	if o.level%3 == 0 && o.orbCount < 5 {
		o.orbCount++
	}
}

func (o *Orbital) Kind() WeaponKind { return KindOrbital }
func (o *Orbital) Stats() WeaponStats {
	return WeaponStats{Name: "Orb", Kind: KindOrbital, Level: o.level, Damage: o.damage, Range: o.radius}
}
func (o *Orbital) AddDamage(v int)             { o.damage += v }
func (o *Orbital) MultiplyCooldown(_ float64)  {} // no cooldown to multiply
func (o *Orbital) AddRange(v float64)          { o.radius += v }
func (o *Orbital) AddPierce(v int) {
	o.orbCount += v
	if o.orbCount > 5 {
		o.orbCount = 5
	}
}
