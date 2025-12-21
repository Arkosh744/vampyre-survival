package weapon

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_Lightning_New(t *testing.T) {
	l := NewLightning()
	require.Equal(t, 1, l.level)
	require.Equal(t, LightningBaseDamage, l.damage)
	require.Equal(t, LightningChainCount, l.chainCount)
}

func Test_Lightning_ChainHitsMultiple(t *testing.T) {
	l := NewLightning()
	owner := physics.Vec2{X: 0, Y: 0}
	targets := []Target{
		{Pos: physics.Vec2{X: 5, Y: 0}, ID: 0},
		{Pos: physics.Vec2{X: 10, Y: 0}, ID: 1},
		{Pos: physics.Vec2{X: 15, Y: 0}, ID: 2},
	}

	l.Update(0.1, owner, targets)
	hits := l.GetHits()
	require.Len(t, hits, 3, "should chain to all 3 targets")

	require.Equal(t, 0, hits[0].TargetID)
	require.Equal(t, 1, hits[1].TargetID)
	require.Equal(t, 2, hits[2].TargetID)
}

func Test_Lightning_DamageDegrades(t *testing.T) {
	l := NewLightning()
	owner := physics.Vec2{X: 0, Y: 0}
	targets := []Target{
		{Pos: physics.Vec2{X: 5, Y: 0}, ID: 0},
		{Pos: physics.Vec2{X: 10, Y: 0}, ID: 1},
		{Pos: physics.Vec2{X: 15, Y: 0}, ID: 2},
	}

	l.Update(0.1, owner, targets)
	hits := l.GetHits()
	require.Len(t, hits, 3)

	require.Equal(t, 25, hits[0].Damage) // base
	require.Equal(t, 17, hits[1].Damage) // 25 * 0.7 = 17.5 → 17
	require.Equal(t, 12, hits[2].Damage) // 17.5 * 0.7 = 12.25 → 12
}

func Test_Lightning_MissesOutOfRange(t *testing.T) {
	l := NewLightning()
	owner := physics.Vec2{X: 0, Y: 0}
	targets := []Target{
		{Pos: physics.Vec2{X: 50, Y: 0}, ID: 0}, // way out of range
	}

	l.Update(0.1, owner, targets)
	require.Empty(t, l.GetHits())
}

func Test_Lightning_CooldownBlocksFire(t *testing.T) {
	l := NewLightning()
	owner := physics.Vec2{X: 0, Y: 0}
	targets := []Target{
		{Pos: physics.Vec2{X: 5, Y: 0}, ID: 0},
	}

	l.Update(0.1, owner, targets)
	require.NotEmpty(t, l.GetHits())

	// Second update with small dt — should be on cooldown
	l.Update(0.1, owner, targets)
	require.Empty(t, l.GetHits())
}

func Test_Lightning_Upgrade(t *testing.T) {
	l := NewLightning()
	l.Upgrade() // level 2: +5 dmg, chainCount++
	require.Equal(t, 2, l.level)
	require.Equal(t, LightningBaseDamage+5, l.damage)
	require.Equal(t, LightningChainCount+1, l.chainCount)
}

func Test_Lightning_AddPierce(t *testing.T) {
	l := NewLightning()
	require.Equal(t, LightningChainCount, l.chainCount)
	l.AddPierce(2)
	require.Equal(t, LightningChainCount+2, l.chainCount)
}

func Test_Lightning_AddPierce_Cap(t *testing.T) {
	l := NewLightning()
	l.AddPierce(10) // should cap at 7
	require.Equal(t, 7, l.chainCount)
}

func Test_Lightning_ChainRadiusLimit(t *testing.T) {
	l := NewLightning()
	owner := physics.Vec2{X: 0, Y: 0}
	// First target close, second target too far from first for chain
	targets := []Target{
		{Pos: physics.Vec2{X: 5, Y: 0}, ID: 0},
		{Pos: physics.Vec2{X: 5 + LightningChainRadius + 5, Y: 0}, ID: 1},
	}

	l.Update(0.1, owner, targets)
	hits := l.GetHits()
	require.Len(t, hits, 1, "second target too far to chain")
}

func Test_Lightning_Visuals(t *testing.T) {
	l := NewLightning()
	owner := physics.Vec2{X: 0, Y: 0}
	targets := []Target{
		{Pos: physics.Vec2{X: 5, Y: 0}, ID: 0},
		{Pos: physics.Vec2{X: 10, Y: 0}, ID: 1},
	}

	l.Update(0.1, owner, targets)
	visuals := l.GetVisuals()
	// 2 targets: each gets 'z' + '~' midpoint = 4 visuals
	require.Len(t, visuals, 4)
}
