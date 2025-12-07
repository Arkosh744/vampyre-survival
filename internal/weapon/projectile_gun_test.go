package weapon

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_ProjectileGun_New(t *testing.T) {
	pg := NewProjectileGun()
	require.Equal(t, 1, pg.Level())
}

func Test_ProjectileGun_FiresAtNearest(t *testing.T) {
	pg := NewProjectileGun()
	pg.cooldownTimer = 0
	owner := physics.Vec2{X: 10, Y: 10}
	targets := []Target{
		{Pos: physics.Vec2{X: 20, Y: 10}, ID: 1},
	}
	pg.Update(0.016, owner, targets)
	require.Equal(t, 1, len(pg.Bullets))
}

func Test_ProjectileGun_BulletMovesToTarget(t *testing.T) {
	pg := NewProjectileGun()
	pg.cooldownTimer = 0
	owner := physics.Vec2{X: 0, Y: 0}
	targets := []Target{{Pos: physics.Vec2{X: 20, Y: 0}, ID: 1}}

	pg.Update(0.016, owner, targets)
	require.Equal(t, 1, len(pg.Bullets))
	startX := pg.Bullets[0].Body.Pos.X

	pg.Update(0.5, owner, targets)
	require.Greater(t, pg.Bullets[0].Body.Pos.X, startX)
}

func Test_ProjectileGun_BulletHitsTarget(t *testing.T) {
	pg := NewProjectileGun()
	pg.Bullets = append(pg.Bullets, Bullet{
		Body:   physics.Body{Pos: physics.Vec2{X: 20, Y: 0}, Width: 1, Height: 1},
		Damage: 10,
		TTL:    5,
	})
	targets := []Target{{Pos: physics.Vec2{X: 20, Y: 0}, ID: 42}}

	pg.checkBulletHits(targets)
	hits := pg.GetHits()
	require.Equal(t, 1, len(hits))
	require.Equal(t, 42, hits[0].TargetID)
}

func Test_ProjectileGun_BulletExpires(t *testing.T) {
	pg := NewProjectileGun()
	pg.Bullets = append(pg.Bullets, Bullet{TTL: 0.01})
	pg.Update(0.1, physics.Vec2{}, nil)
	require.Equal(t, 0, len(pg.Bullets))
}
