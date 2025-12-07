package weapon

import "github.com/arkosh/vampyre-survival/internal/physics"

type Target struct {
	Pos physics.Vec2
	ID  int
}

type HitResult struct {
	TargetID int
	Damage   int
}

type WeaponKind string

const (
	KindSword      WeaponKind = "sword"
	KindProjectile WeaponKind = "projectile"
	KindAoE        WeaponKind = "aoe"
)

type Weapon interface {
	Update(dt float64, ownerPos physics.Vec2, targets []Target)
	GetHits() []HitResult
	GetVisuals() []Visual
	Level() int
	Upgrade()
	Kind() WeaponKind
	AddDamage(v int)
	MultiplyCooldown(factor float64)
	AddRange(v float64)
}

type Visual struct {
	Pos  physics.Vec2
	Char rune
	TTL  float64
}
