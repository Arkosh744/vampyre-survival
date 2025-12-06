package engine

import (
	"math/rand"

	"github.com/arkosh/vampyre-survival/internal/physics"
)

type Camera struct {
	Pos         physics.Vec2
	Width       int
	Height      int
	WorldWidth  float64
	WorldHeight float64

	ShakeOffset physics.Vec2
	ShakeTimer  float64
}

func NewCamera(width, height int) *Camera {
	return &Camera{Width: width, Height: height}
}

func (c *Camera) SetWorldBounds(w, h float64) {
	c.WorldWidth = w
	c.WorldHeight = h
}

func (c *Camera) Follow(target physics.Vec2, dt float64) {
	c.Pos.X = target.X - float64(c.Width)/2
	c.Pos.Y = target.Y - float64(c.Height)/2

	// Decay shake
	if c.ShakeTimer > 0 {
		c.ShakeTimer -= dt
		if c.ShakeTimer <= 0 {
			c.ShakeTimer = 0
			c.ShakeOffset = physics.Vec2{}
		}
	}

	c.clamp()
}

func (c *Camera) Shake(intensity float64) {
	c.ShakeOffset = physics.Vec2{
		X: (rand.Float64()*2 - 1) * intensity,
		Y: (rand.Float64()*2 - 1) * intensity,
	}
	c.ShakeTimer = 0.15
}

func (c *Camera) clamp() {
	if c.Pos.X < 0 {
		c.Pos.X = 0
	}
	if c.Pos.Y < 0 {
		c.Pos.Y = 0
	}
	if c.WorldWidth > 0 && c.Pos.X+float64(c.Width) > c.WorldWidth {
		c.Pos.X = c.WorldWidth - float64(c.Width)
	}
	if c.WorldHeight > 0 && c.Pos.Y+float64(c.Height) > c.WorldHeight {
		c.Pos.Y = c.WorldHeight - float64(c.Height)
	}
}

func (c *Camera) WorldToScreen(worldPos physics.Vec2) (int, int) {
	return int(worldPos.X - c.Pos.X + c.ShakeOffset.X), int(worldPos.Y - c.Pos.Y + c.ShakeOffset.Y)
}

func (c *Camera) IsVisible(worldPos physics.Vec2) bool {
	sx, sy := c.WorldToScreen(worldPos)
	return sx >= -1 && sx <= c.Width && sy >= -1 && sy <= c.Height
}
