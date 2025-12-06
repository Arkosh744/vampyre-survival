package engine

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_Camera_New(t *testing.T) {
	c := NewCamera(80, 24)
	require.Equal(t, 80, c.Width)
	require.Equal(t, 24, c.Height)
}

func Test_Camera_Follow(t *testing.T) {
	c := NewCamera(80, 24)
	target := physics.Vec2{X: 100, Y: 100}
	c.Follow(target, 0.033)
	require.InDelta(t, 60.0, c.Pos.X, 0.001)
	require.InDelta(t, 88.0, c.Pos.Y, 0.001)
}

func Test_Camera_WorldToScreen(t *testing.T) {
	c := NewCamera(80, 24)
	c.Pos = physics.Vec2{X: 10, Y: 5}
	sx, sy := c.WorldToScreen(physics.Vec2{X: 15, Y: 8})
	require.Equal(t, 5, sx)
	require.Equal(t, 3, sy)
}

func Test_Camera_IsVisible(t *testing.T) {
	c := NewCamera(80, 24)
	c.Pos = physics.Vec2{X: 0, Y: 0}
	require.True(t, c.IsVisible(physics.Vec2{X: 40, Y: 12}))
	require.False(t, c.IsVisible(physics.Vec2{X: 200, Y: 200}))
}

func Test_Camera_Clamp(t *testing.T) {
	c := NewCamera(80, 24)
	c.SetWorldBounds(200, 200)
	c.Follow(physics.Vec2{X: 0, Y: 0}, 0.033)
	require.GreaterOrEqual(t, c.Pos.X, 0.0)
	require.GreaterOrEqual(t, c.Pos.Y, 0.0)
}

func Test_Camera_ShakeDecay(t *testing.T) {
	c := NewCamera(80, 24)
	c.Shake(2.0)
	require.Greater(t, c.ShakeTimer, 0.0)
	require.True(t, c.ShakeOffset.X != 0 || c.ShakeOffset.Y != 0)

	// After enough time, shake decays
	c.Follow(physics.Vec2{X: 40, Y: 12}, 0.2)
	require.Equal(t, 0.0, c.ShakeTimer)
	require.InDelta(t, 0.0, c.ShakeOffset.X, 0.001)
	require.InDelta(t, 0.0, c.ShakeOffset.Y, 0.001)
}
