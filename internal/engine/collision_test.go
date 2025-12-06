package engine

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_AABB_Overlap(t *testing.T) {
	a := physics.Body{Pos: physics.Vec2{X: 0, Y: 0}, Width: 2, Height: 2}
	b := physics.Body{Pos: physics.Vec2{X: 1, Y: 1}, Width: 2, Height: 2}
	require.True(t, CheckAABB(&a, &b))
}

func Test_AABB_NoOverlap(t *testing.T) {
	a := physics.Body{Pos: physics.Vec2{X: 0, Y: 0}, Width: 2, Height: 2}
	b := physics.Body{Pos: physics.Vec2{X: 5, Y: 5}, Width: 2, Height: 2}
	require.False(t, CheckAABB(&a, &b))
}

func Test_AABB_TouchingEdge(t *testing.T) {
	a := physics.Body{Pos: physics.Vec2{X: 0, Y: 0}, Width: 2, Height: 2}
	b := physics.Body{Pos: physics.Vec2{X: 2, Y: 0}, Width: 2, Height: 2}
	require.False(t, CheckAABB(&a, &b))
}

func Test_CircleCollision_Overlap(t *testing.T) {
	a := physics.Vec2{X: 0, Y: 0}
	b := physics.Vec2{X: 2, Y: 0}
	require.True(t, CheckCircle(a, b, 1.5, 1.5))
}

func Test_CircleCollision_NoOverlap(t *testing.T) {
	a := physics.Vec2{X: 0, Y: 0}
	b := physics.Vec2{X: 10, Y: 0}
	require.False(t, CheckCircle(a, b, 1.0, 1.0))
}

func Test_PointInRect(t *testing.T) {
	p := physics.Vec2{X: 5, Y: 5}
	require.True(t, PointInRect(p, physics.Vec2{X: 0, Y: 0}, 10, 10))
	require.False(t, PointInRect(p, physics.Vec2{X: 6, Y: 6}, 2, 2))
}
