package game

import "fmt"

const (
	ComboTimeout       = 2.0
	ComboLabelDuration = 1.5
)

type Combo struct {
	Count      int
	Timer      float64
	BestStreak int
	Frozen     bool // Momentum reward: combo never expires
}

func NewCombo() *Combo {
	return &Combo{}
}

func (c *Combo) RegisterKill() {
	c.Timer = 0
	c.Count++
	if c.Count > c.BestStreak {
		c.BestStreak = c.Count
	}
}

func (c *Combo) Update(dt float64) {
	if c.Count == 0 {
		return
	}
	if c.Frozen {
		return
	}
	c.Timer += dt
	if c.Timer >= ComboTimeout {
		c.Count = 0
		c.Timer = 0
	}
}

// XPMult returns XP multiplier based on current combo tier.
func (c *Combo) XPMult() float64 {
	switch {
	case c.Count >= 20:
		return 3.0
	case c.Count >= 10:
		return 2.0
	case c.Count >= 5:
		return 1.5
	default:
		return 1.0
	}
}

// DamageMult returns damage multiplier: linear +0.005% per kill.
func (c *Combo) DamageMult() float64 {
	return 1.0 + float64(c.Count)*0.00005
}

// Label returns combo banner text, empty string if combo < 5.
func (c *Combo) Label() string {
	switch {
	case c.Count >= 25:
		return fmt.Sprintf("UNSTOPPABLE! x%d", c.Count)
	case c.Count >= 15:
		return fmt.Sprintf("MEGA KILL! x%d", c.Count)
	case c.Count >= 5:
		return fmt.Sprintf("COMBO x%d!", c.Count)
	default:
		return ""
	}
}
