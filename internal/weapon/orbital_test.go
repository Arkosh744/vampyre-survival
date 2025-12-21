package weapon

import (
	"math"
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_Orbital_New(t *testing.T) {
	o := NewOrbital()
	require.Equal(t, 1, o.level)
	require.Equal(t, OrbitalBaseDamage, o.damage)
	require.Equal(t, OrbitalBaseOrbCount, o.orbCount)
}

func Test_Orbital_OrbPositions(t *testing.T) {
	o := NewOrbital()
	owner := physics.Vec2{X: 50, Y: 50}
	// angle=0: orb0 at (50+4, 50), orb1 at (50-4, 50) — opposite sides
	o.Update(0, owner, nil)

	visuals := o.GetVisuals()
	require.Len(t, visuals, 2)

	// Orb 0: angle=0 → cos(0)=1, sin(0)=0
	require.InDelta(t, 54.0, visuals[0].Pos.X, 0.01)
	require.InDelta(t, 50.0, visuals[0].Pos.Y, 0.01)

	// Orb 1: angle=π → cos(π)=-1, sin(π)=0
	require.InDelta(t, 46.0, visuals[1].Pos.X, 0.01)
	require.InDelta(t, 50.0, visuals[1].Pos.Y, 0.01)
}

func Test_Orbital_HitsNearbyTarget(t *testing.T) {
	o := NewOrbital()
	owner := physics.Vec2{X: 50, Y: 50}
	// Place target where orb 0 will be (radius=4 right of owner)
	targets := []Target{
		{Pos: physics.Vec2{X: 54, Y: 50}, ID: 0},
	}

	o.Update(0.01, owner, targets)
	hits := o.GetHits()
	require.Len(t, hits, 1)
	require.Equal(t, 0, hits[0].TargetID)
	require.Equal(t, OrbitalBaseDamage, hits[0].Damage)
}

func Test_Orbital_PerTargetCooldown(t *testing.T) {
	o := NewOrbital()
	owner := physics.Vec2{X: 50, Y: 50}

	// Place target at orb0 initial position (angle=0, radius=4 → X=54)
	targets := []Target{
		{Pos: physics.Vec2{X: 54, Y: 50}, ID: 0},
	}

	// First hit
	o.Update(0.01, owner, targets)
	require.Len(t, o.GetHits(), 1)

	// Second update within cooldown — no hit even if orb is near target
	// Use dt=0 to prevent orb rotation away from target
	o.Update(0, owner, targets)
	require.Empty(t, o.GetHits())

	// Advance past cooldown but rotate back to same position (full rotation = 2π/speed)
	// We expire the cooldown manually by updating with enough time, then
	// rotate angle to bring orb back to target.
	fullRotation := 2 * math.Pi / OrbitalRotationSpeed
	o.Update(fullRotation, owner, targets) // orb back at start, cooldown expired
	require.Len(t, o.GetHits(), 1)
}

func Test_Orbital_MissesDistant(t *testing.T) {
	o := NewOrbital()
	owner := physics.Vec2{X: 50, Y: 50}
	targets := []Target{
		{Pos: physics.Vec2{X: 100, Y: 100}, ID: 0},
	}

	o.Update(0.01, owner, targets)
	require.Empty(t, o.GetHits())
}

func Test_Orbital_UpgradeAddsOrb(t *testing.T) {
	o := NewOrbital()
	require.Equal(t, 2, o.orbCount)

	o.Upgrade() // level 2
	require.Equal(t, 2, o.orbCount)

	o.Upgrade() // level 3 — adds orb
	require.Equal(t, 3, o.orbCount)
	require.Equal(t, OrbitalBaseDamage+6, o.damage)
}

func Test_Orbital_AddPierce_AddsOrb(t *testing.T) {
	o := NewOrbital()
	o.AddPierce(1)
	require.Equal(t, OrbitalBaseOrbCount+1, o.orbCount)
}

func Test_Orbital_AddPierce_Cap(t *testing.T) {
	o := NewOrbital()
	o.AddPierce(10)
	require.Equal(t, 5, o.orbCount)
}

func Test_Orbital_Rotation(t *testing.T) {
	o := NewOrbital()
	owner := physics.Vec2{X: 50, Y: 50}

	// After π/2 rotation (quarter turn), orb 0 should be above owner
	dt := (math.Pi / 2) / OrbitalRotationSpeed
	o.Update(dt, owner, nil)

	visuals := o.GetVisuals()
	require.Len(t, visuals, 2)

	// Orb 0: angle=π/2 → cos=0, sin=1
	require.InDelta(t, 50.0, visuals[0].Pos.X, 0.1)
	require.InDelta(t, 54.0, visuals[0].Pos.Y, 0.1)
}
