package physics

type Body struct {
	Pos      Vec2
	Vel      Vec2
	Friction float64
	MaxSpeed float64
	Width    float64
	Height   float64
}

func (b *Body) Update(dt float64) {
	b.Pos = b.Pos.Add(b.Vel.Scale(dt))

	if b.Friction > 0 {
		b.Vel = b.Vel.Scale(b.Friction)
	}

	if b.MaxSpeed > 0 && b.Vel.Length() > b.MaxSpeed {
		b.Vel = b.Vel.Normalize().Scale(b.MaxSpeed)
	}
}

func (b *Body) ApplyForce(force Vec2) {
	b.Vel = b.Vel.Add(force)
}
