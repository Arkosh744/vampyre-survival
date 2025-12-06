package entity

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_Pickup_NewXPGem(t *testing.T) {
	p := NewXPGem(5, 10, 1)
	require.Equal(t, 5.0, p.Body.Pos.X)
	require.Equal(t, 10.0, p.Body.Pos.Y)
	require.Equal(t, 1, p.Value)
	require.Equal(t, PickupXP, p.Type)
}

func Test_Pickup_MagnetToward(t *testing.T) {
	p := NewXPGem(20, 0, 1)
	playerPos := physics.Vec2{X: 18, Y: 0}
	p.MagnetToward(playerPos, 3.0)
	p.Update(0.5)
	require.Less(t, p.Body.Pos.X, 20.0)
}

func Test_Pickup_NoMagnetWhenFar(t *testing.T) {
	p := NewXPGem(100, 0, 1)
	playerPos := physics.Vec2{X: 0, Y: 0}
	p.MagnetToward(playerPos, 3.0)
	require.InDelta(t, 0.0, p.Body.Vel.Length(), 0.001)
}

func Test_Pickup_Collected(t *testing.T) {
	p := NewXPGem(10, 10, 5)
	require.False(t, p.Collected)
	p.Collect()
	require.True(t, p.Collected)
}
