package weapon

import "github.com/arkosh/vampyre-survival/internal/physics"

const (
	LightningBaseDamage  = 25
	LightningBaseCooldown = 0.6
	LightningBaseRange   = 12.0
	LightningChainCount  = 3
	LightningChainRadius = 8.0
	LightningDegradeRate = 0.7
)

type Lightning struct {
	level         int
	cooldownTimer float64
	cooldown      float64
	damage        int
	attackRange   float64
	chainCount    int
	hits          []HitResult
	visuals       []Visual
}

func NewLightning() *Lightning {
	return &Lightning{
		level:       1,
		cooldown:    LightningBaseCooldown,
		damage:      LightningBaseDamage,
		attackRange: LightningBaseRange,
		chainCount:  LightningChainCount,
	}
}

func (l *Lightning) Update(dt float64, ownerPos physics.Vec2, targets []Target) {
	l.hits = nil
	l.visuals = nil

	if l.cooldownTimer > 0 {
		l.cooldownTimer -= dt
		return
	}

	if len(targets) == 0 {
		return
	}

	// Find nearest target within attack range
	firstIdx := -1
	bestDist := l.attackRange + 1
	for i, t := range targets {
		d := ownerPos.DistanceTo(t.Pos)
		if d < bestDist {
			bestDist = d
			firstIdx = i
		}
	}
	if firstIdx < 0 {
		return
	}

	l.cooldownTimer = l.cooldown

	hit := make(map[int]bool)
	dmg := float64(l.damage)
	prevPos := ownerPos
	curIdx := firstIdx

	for hop := 0; hop < l.chainCount; hop++ {
		t := targets[curIdx]
		hit[curIdx] = true

		intDmg := int(dmg)
		if intDmg < 1 {
			intDmg = 1
		}
		l.hits = append(l.hits, HitResult{TargetID: t.ID, Damage: intDmg})

		// Visual: 'z' on target
		l.visuals = append(l.visuals, Visual{Pos: t.Pos, Char: 'z', TTL: 0.25})

		// Visual: '~' at midpoint from previous position
		mid := physics.Vec2{
			X: (prevPos.X + t.Pos.X) / 2,
			Y: (prevPos.Y + t.Pos.Y) / 2,
		}
		l.visuals = append(l.visuals, Visual{Pos: mid, Char: '~', TTL: 0.2})

		// Find next chain target
		prevPos = t.Pos
		dmg *= LightningDegradeRate

		nextIdx := -1
		nextDist := LightningChainRadius + 1
		for i, nt := range targets {
			if hit[i] {
				continue
			}
			d := t.Pos.DistanceTo(nt.Pos)
			if d < nextDist {
				nextDist = d
				nextIdx = i
			}
		}
		if nextIdx < 0 {
			break
		}
		curIdx = nextIdx
	}
}

func (l *Lightning) GetHits() []HitResult { return l.hits }
func (l *Lightning) GetVisuals() []Visual  { return l.visuals }
func (l *Lightning) Level() int            { return l.level }

func (l *Lightning) Upgrade() {
	l.level++
	l.damage += 5
	l.cooldown *= 0.9
	if l.level%2 == 0 && l.chainCount < 7 {
		l.chainCount++
	}
}

func (l *Lightning) Kind() WeaponKind            { return KindLightning }
func (l *Lightning) AddDamage(v int)             { l.damage += v }
func (l *Lightning) MultiplyCooldown(f float64)  { l.cooldown *= f }
func (l *Lightning) AddRange(v float64)          { l.attackRange += v }
func (l *Lightning) AddPierce(v int) {
	l.chainCount += v
	if l.chainCount > 7 {
		l.chainCount = 7
	}
}
