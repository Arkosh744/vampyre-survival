package weapon

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_Sword_New(t *testing.T) {
	s := NewSword()
	require.Equal(t, 1, s.Level())
	require.Equal(t, 0, len(s.GetHits()))
	require.Equal(t, KindSword, s.Kind())
}

func Test_Sword_HitsNearbyEnemy(t *testing.T) {
	s := NewSword()
	s.cooldownTimer = 0
	ownerPos := physics.Vec2{X: 10, Y: 10}
	targets := []Target{
		{Pos: physics.Vec2{X: 12, Y: 10}, ID: 1},
	}
	s.Update(0.016, ownerPos, targets)
	hits := s.GetHits()
	require.Equal(t, 1, len(hits))
	require.Equal(t, 1, hits[0].TargetID)
}

func Test_Sword_MissesDistantEnemy(t *testing.T) {
	s := NewSword()
	s.cooldownTimer = 0
	ownerPos := physics.Vec2{X: 10, Y: 10}
	targets := []Target{
		{Pos: physics.Vec2{X: 50, Y: 50}, ID: 1},
	}
	s.Update(0.016, ownerPos, targets)
	hits := s.GetHits()
	require.Equal(t, 0, len(hits))
}

func Test_Sword_Cooldown(t *testing.T) {
	s := NewSword()
	s.cooldownTimer = 0
	ownerPos := physics.Vec2{X: 10, Y: 10}
	targets := []Target{{Pos: physics.Vec2{X: 12, Y: 10}, ID: 1}}

	s.Update(0.016, ownerPos, targets)
	require.Equal(t, 1, len(s.GetHits()))

	s.Update(0.016, ownerPos, targets)
	require.Equal(t, 0, len(s.GetHits()))
}

func Test_Sword_Upgrade(t *testing.T) {
	s := NewSword()
	s.Upgrade()
	require.Equal(t, 2, s.Level())
}

func Test_Sword_AddDamage(t *testing.T) {
	s := NewSword()
	oldDmg := s.damage
	s.AddDamage(10)
	require.Equal(t, oldDmg+10, s.damage)
}

func Test_Sword_MultiplyCooldown(t *testing.T) {
	s := NewSword()
	oldCD := s.cooldown
	s.MultiplyCooldown(0.5)
	require.InDelta(t, oldCD*0.5, s.cooldown, 0.001)
}

func Test_Sword_AddRange(t *testing.T) {
	s := NewSword()
	oldRange := s.attackRange
	s.AddRange(3.0)
	require.InDelta(t, oldRange+3.0, s.attackRange, 0.001)
}
