package physics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Vec2_Add(t *testing.T) {
	a := Vec2{X: 1, Y: 2}
	b := Vec2{X: 3, Y: 4}
	result := a.Add(b)
	require.Equal(t, Vec2{X: 4, Y: 6}, result)
}

func Test_Vec2_Sub(t *testing.T) {
	a := Vec2{X: 5, Y: 7}
	b := Vec2{X: 2, Y: 3}
	result := a.Sub(b)
	require.Equal(t, Vec2{X: 3, Y: 4}, result)
}

func Test_Vec2_Scale(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	result := v.Scale(2)
	require.Equal(t, Vec2{X: 6, Y: 8}, result)
}

func Test_Vec2_Length(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	require.InDelta(t, 5.0, v.Length(), 0.001)
}

func Test_Vec2_Normalize(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	n := v.Normalize()
	require.InDelta(t, 1.0, n.Length(), 0.001)
	require.InDelta(t, 0.6, n.X, 0.001)
	require.InDelta(t, 0.8, n.Y, 0.001)
}

func Test_Vec2_Normalize_Zero(t *testing.T) {
	v := Vec2{X: 0, Y: 0}
	n := v.Normalize()
	require.Equal(t, Vec2{X: 0, Y: 0}, n)
}

func Test_Vec2_DistanceTo(t *testing.T) {
	a := Vec2{X: 0, Y: 0}
	b := Vec2{X: 3, Y: 4}
	require.InDelta(t, 5.0, a.DistanceTo(b), 0.001)
}

func Test_Vec2_Dot(t *testing.T) {
	a := Vec2{X: 1, Y: 0}
	b := Vec2{X: 0, Y: 1}
	require.InDelta(t, 0.0, a.Dot(b), 0.001)
}
