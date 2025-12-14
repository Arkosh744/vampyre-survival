package weapon

import "github.com/arkosh/vampyre-survival/internal/physics"

const (
	SwordBaseRange    = 5.0
	SwordBaseDamage   = 25
	SwordBaseCooldown = 0.3
)

type Sword struct {
	level         int
	cooldownTimer float64
	cooldown      float64
	damage        int
	attackRange   float64
	hits          []HitResult
	visuals       []Visual
}

func NewSword() *Sword {
	return &Sword{
		level:       1,
		cooldown:    SwordBaseCooldown,
		damage:      SwordBaseDamage,
		attackRange: SwordBaseRange,
	}
}

func (s *Sword) Update(dt float64, ownerPos physics.Vec2, targets []Target) {
	s.hits = nil
	s.visuals = nil

	if s.cooldownTimer > 0 {
		s.cooldownTimer -= dt
		return
	}

	var nearest *Target
	nearestDist := s.attackRange + 1

	for i := range targets {
		d := ownerPos.DistanceTo(targets[i].Pos)
		if d < nearestDist {
			nearestDist = d
			nearest = &targets[i]
		}
	}

	if nearest == nil {
		return
	}

	for i := range targets {
		d := ownerPos.DistanceTo(targets[i].Pos)
		if d <= s.attackRange {
			s.hits = append(s.hits, HitResult{
				TargetID: targets[i].ID,
				Damage:   s.damage,
			})
		}
	}

	if len(s.hits) > 0 {
		s.cooldownTimer = s.cooldown
		dir := nearest.Pos.Sub(ownerPos).Normalize()
		s.visuals = append(s.visuals, Visual{
			Pos:  ownerPos.Add(dir),
			Char: '╱',
			TTL:  0.2,
		})
	}
}

func (s *Sword) GetHits() []HitResult { return s.hits }
func (s *Sword) GetVisuals() []Visual  { return s.visuals }
func (s *Sword) Level() int            { return s.level }

func (s *Sword) Upgrade() {
	s.level++
	s.damage += 5
	s.cooldown *= 0.85
	s.attackRange += 0.5
}

func (s *Sword) Kind() WeaponKind            { return KindSword }
func (s *Sword) AddDamage(v int)             { s.damage += v }
func (s *Sword) MultiplyCooldown(f float64)  { s.cooldown *= f }
func (s *Sword) AddRange(v float64)          { s.attackRange += v }
func (s *Sword) AddPierce(_ int)             {}
