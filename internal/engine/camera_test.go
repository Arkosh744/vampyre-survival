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
	// Smoothness=0 (default) → instant snap, same behavior as before.
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

func Test_Camera_WorldToScreenPx(t *testing.T) {
	c := NewCamera(80, 24)
	c.Pos = physics.Vec2{X: 10, Y: 5}
	px, py := c.WorldToScreenPx(physics.Vec2{X: 15, Y: 8})
	require.InDelta(t, 5.0*TileSize, px, 0.001)
	require.InDelta(t, 3.0*TileSize, py, 0.001)
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

func Test_Camera_SmoothFollow(t *testing.T) {
	c := NewCamera(80, 24)
	c.Smoothness = 0.1

	// Start camera at origin.
	c.Pos = physics.Vec2{X: 0, Y: 0}

	target := physics.Vec2{X: 100, Y: 100}
	desired := physics.Vec2{
		X: target.X - float64(c.Width)/2,
		Y: target.Y - float64(c.Height)/2,
	}

	// One frame of lerp should move toward desired but not reach it.
	c.Follow(target, 0.033)

	require.Greater(t, c.Pos.X, 0.0, "camera should move toward target X")
	require.Less(t, c.Pos.X, desired.X, "camera should not reach target X in one frame")
	require.Greater(t, c.Pos.Y, 0.0, "camera should move toward target Y")
	require.Less(t, c.Pos.Y, desired.Y, "camera should not reach target Y in one frame")

	// After many frames the camera should converge close to desired.
	for i := 0; i < 300; i++ {
		c.Follow(target, 0.016)
	}
	require.InDelta(t, desired.X, c.Pos.X, 0.01)
	require.InDelta(t, desired.Y, c.Pos.Y, 0.01)
}
