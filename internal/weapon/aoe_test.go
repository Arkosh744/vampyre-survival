package weapon

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_AoE_New(t *testing.T) {
	a := NewAoE()
	require.Equal(t, 1, a.Level())
}

func Test_AoE_HitsNearbyEnemies(t *testing.T) {
	a := NewAoE()
	a.cooldownTimer = 0
	owner := physics.Vec2{X: 10, Y: 10}
	targets := []Target{
		{Pos: physics.Vec2{X: 12, Y: 10}, ID: 1},
		{Pos: physics.Vec2{X: 14, Y: 10}, ID: 2},
	}
	a.Update(0.016, owner, targets)
	hits := a.GetHits()
	require.Equal(t, 2, len(hits))
}

func Test_AoE_MissesDistantEnemies(t *testing.T) {
	a := NewAoE()
	a.cooldownTimer = 0
	owner := physics.Vec2{X: 10, Y: 10}
	targets := []Target{
		{Pos: physics.Vec2{X: 50, Y: 50}, ID: 1},
	}
	a.Update(0.016, owner, targets)
	require.Equal(t, 0, len(a.GetHits()))
}

func Test_AoE_Cooldown(t *testing.T) {
	a := NewAoE()
	a.cooldownTimer = 0
	owner := physics.Vec2{X: 10, Y: 10}
	targets := []Target{{Pos: physics.Vec2{X: 12, Y: 10}, ID: 1}}

	a.Update(0.016, owner, targets)
	require.Equal(t, 1, len(a.GetHits()))

	a.Update(0.016, owner, targets)
	require.Equal(t, 0, len(a.GetHits()))
}

func Test_AoE_Upgrade(t *testing.T) {
	a := NewAoE()
	r1 := a.radius
	a.Upgrade()
	require.Equal(t, 2, a.Level())
	require.Greater(t, a.radius, r1)
}
