package physics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Body_Update_AppliesVelocity(t *testing.T) {
	b := Body{
		Pos: Vec2{X: 10, Y: 10},
		Vel: Vec2{X: 5, Y: -3},
	}
	b.Update(1.0)
	require.InDelta(t, 15.0, b.Pos.X, 0.001)
	require.InDelta(t, 7.0, b.Pos.Y, 0.001)
}

func Test_Body_Update_AppliesFriction(t *testing.T) {
	b := Body{
		Pos:      Vec2{X: 0, Y: 0},
		Vel:      Vec2{X: 10, Y: 10},
		Friction: 0.9,
	}
	b.Update(1.0)
	require.InDelta(t, 9.0, b.Vel.X, 0.001)
	require.InDelta(t, 9.0, b.Vel.Y, 0.001)
}

func Test_Body_ApplyForce(t *testing.T) {
	b := Body{Vel: Vec2{X: 0, Y: 0}}
	b.ApplyForce(Vec2{X: 5, Y: 3})
	require.InDelta(t, 5.0, b.Vel.X, 0.001)
	require.InDelta(t, 3.0, b.Vel.Y, 0.001)
}

func Test_Body_ClampSpeed(t *testing.T) {
	b := Body{
		Vel:      Vec2{X: 100, Y: 0},
		MaxSpeed: 10,
	}
	b.Update(0.016)
	require.LessOrEqual(t, b.Vel.Length(), 10.0+0.001)
}
