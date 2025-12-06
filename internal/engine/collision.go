package engine

import "github.com/arkosh/vampyre-survival/internal/physics"

func CheckAABB(a, b *physics.Body) bool {
	return a.Pos.X < b.Pos.X+b.Width &&
		a.Pos.X+a.Width > b.Pos.X &&
		a.Pos.Y < b.Pos.Y+b.Height &&
		a.Pos.Y+a.Height > b.Pos.Y
}

func CheckCircle(a, b physics.Vec2, radiusA, radiusB float64) bool {
	dist := a.DistanceTo(b)
	return dist < radiusA+radiusB
}

func PointInRect(p, rectPos physics.Vec2, w, h float64) bool {
	return p.X >= rectPos.X && p.X <= rectPos.X+w &&
		p.Y >= rectPos.Y && p.Y <= rectPos.Y+h
}
