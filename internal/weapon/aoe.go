package weapon

import "github.com/arkosh/vampyre-survival/internal/physics"

const (
	AoEBaseRadius   = 8.0
	AoEBaseDamage   = 15
	AoEBaseCooldown = 0.8
)

type AoE struct {
	level         int
	cooldownTimer float64
	cooldown      float64
	damage        int
	radius        float64
	hits          []HitResult
	visuals       []Visual
}

func NewAoE() *AoE {
	return &AoE{
		level:    1,
		cooldown: AoEBaseCooldown,
		damage:   AoEBaseDamage,
		radius:   AoEBaseRadius,
	}
}

func (a *AoE) Update(dt float64, ownerPos physics.Vec2, targets []Target) {
	a.hits = nil
	a.visuals = nil

	if a.cooldownTimer > 0 {
		a.cooldownTimer -= dt
		return
	}

	fired := false
	for _, t := range targets {
		if ownerPos.DistanceTo(t.Pos) <= a.radius {
			a.hits = append(a.hits, HitResult{TargetID: t.ID, Damage: a.damage})
			fired = true
		}
	}

	if fired {
		a.cooldownTimer = a.cooldown
		a.visuals = append(a.visuals, Visual{Pos: ownerPos, Char: '◎', TTL: 0.3})
	}
}

func (a *AoE) GetHits() []HitResult { return a.hits }
func (a *AoE) GetVisuals() []Visual  { return a.visuals }
func (a *AoE) Level() int            { return a.level }
func (a *AoE) Upgrade() {
	a.level++
	a.damage += 3
	a.radius += 1.0
	a.cooldown *= 0.9
}

func (a *AoE) Kind() WeaponKind            { return KindAoE }
func (a *AoE) AddDamage(v int)             { a.damage += v }
func (a *AoE) MultiplyCooldown(f float64)  { a.cooldown *= f }
func (a *AoE) AddRange(v float64)          { a.radius += v }
