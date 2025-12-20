package weapon

import "github.com/arkosh/vampyre-survival/internal/physics"

const (
	ProjBaseDamage   = 15
	ProjBaseCooldown = 0.25
	ProjSpeed        = 30.0
	ProjTTL          = 0.85
	ProjHitRadius    = 1.5
)

type Bullet struct {
	Body   physics.Body
	Dir    physics.Vec2
	Damage int
	TTL    float64
	Pierce int
}

type ProjectileGun struct {
	level         int
	cooldownTimer float64
	cooldown      float64
	damage        int
	pierce        int
	Bullets       []Bullet
	hits          []HitResult
	visuals       []Visual
}

func NewProjectileGun() *ProjectileGun {
	return &ProjectileGun{
		level:    1,
		cooldown: ProjBaseCooldown,
		damage:   ProjBaseDamage,
	}
}

func (pg *ProjectileGun) Update(dt float64, ownerPos physics.Vec2, targets []Target) {
	pg.hits = nil
	pg.visuals = nil

	for i := len(pg.Bullets) - 1; i >= 0; i-- {
		pg.Bullets[i].TTL -= dt
		pg.Bullets[i].Body.Pos = pg.Bullets[i].Body.Pos.Add(pg.Bullets[i].Dir.Scale(ProjSpeed * dt))
		if pg.Bullets[i].TTL <= 0 {
			pg.Bullets = append(pg.Bullets[:i], pg.Bullets[i+1:]...)
		}
	}

	pg.checkBulletHits(targets)

	if pg.cooldownTimer > 0 {
		pg.cooldownTimer -= dt
		return
	}

	if len(targets) == 0 {
		return
	}

	nearest := targets[0]
	nearestDist := ownerPos.DistanceTo(nearest.Pos)
	for _, t := range targets[1:] {
		d := ownerPos.DistanceTo(t.Pos)
		if d < nearestDist {
			nearestDist = d
			nearest = t
		}
	}

	dir := nearest.Pos.Sub(ownerPos).Normalize()
	pg.Bullets = append(pg.Bullets, Bullet{
		Body:   physics.Body{Pos: ownerPos.Add(dir), Width: 1, Height: 1},
		Dir:    dir,
		Damage: pg.damage,
		TTL:    ProjTTL,
		Pierce: pg.pierce,
	})
	pg.cooldownTimer = pg.cooldown
}

func (pg *ProjectileGun) checkBulletHits(targets []Target) {
	for i := len(pg.Bullets) - 1; i >= 0; i-- {
		for _, t := range targets {
			if pg.Bullets[i].Body.Pos.DistanceTo(t.Pos) < ProjHitRadius {
				pg.hits = append(pg.hits, HitResult{TargetID: t.ID, Damage: pg.Bullets[i].Damage})
				if pg.Bullets[i].Pierce > 0 {
					pg.Bullets[i].Pierce--
				} else {
					pg.Bullets = append(pg.Bullets[:i], pg.Bullets[i+1:]...)
				}
				break
			}
		}
	}
}

func (pg *ProjectileGun) GetHits() []HitResult { return pg.hits }
func (pg *ProjectileGun) GetVisuals() []Visual {
	for _, b := range pg.Bullets {
		pg.visuals = append(pg.visuals, Visual{Pos: b.Body.Pos, Char: '•', TTL: 0})
	}
	return pg.visuals
}
func (pg *ProjectileGun) Level() int { return pg.level }
func (pg *ProjectileGun) Upgrade() {
	pg.level++
	pg.damage += 4
	pg.cooldown *= 0.88
}

func (pg *ProjectileGun) Kind() WeaponKind            { return KindProjectile }
func (pg *ProjectileGun) AddDamage(v int)             { pg.damage += v }
func (pg *ProjectileGun) MultiplyCooldown(f float64)  { pg.cooldown *= f }
func (pg *ProjectileGun) AddRange(v float64)          { _ = v } // range is TTL-based, no direct range field
func (pg *ProjectileGun) AddPierce(v int)              { pg.pierce += v }
