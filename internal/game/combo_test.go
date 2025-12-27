package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Combo_New(t *testing.T) {
	c := NewCombo()
	require.Equal(t, 0, c.Count)
	require.Equal(t, 0, c.BestStreak)
	require.Equal(t, 0.0, c.Timer)
}

func Test_Combo_RegisterKill_Increments(t *testing.T) {
	c := NewCombo()
	c.RegisterKill()
	require.Equal(t, 1, c.Count)
	c.RegisterKill()
	require.Equal(t, 2, c.Count)
}

func Test_Combo_RegisterKill_ResetsTimer(t *testing.T) {
	c := NewCombo()
	c.RegisterKill()
	c.Update(1.0)
	require.Equal(t, 1.0, c.Timer)
	c.RegisterKill()
	require.Equal(t, 0.0, c.Timer)
}

func Test_Combo_BestStreak(t *testing.T) {
	c := NewCombo()
	for i := 0; i < 10; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 10, c.BestStreak)

	// Reset via timeout
	c.Update(3.0)
	require.Equal(t, 0, c.Count)
	require.Equal(t, 10, c.BestStreak)

	// New streak of 5 doesn't beat old record
	for i := 0; i < 5; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 10, c.BestStreak)
}

func Test_Combo_Timeout_ResetsToZero(t *testing.T) {
	c := NewCombo()
	for i := 0; i < 5; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 5, c.Count)

	c.Update(ComboTimeout)
	require.Equal(t, 0, c.Count)
}

func Test_Combo_NoTimeout_WhileKilling(t *testing.T) {
	c := NewCombo()
	c.RegisterKill()
	c.Update(1.9) // almost timeout
	require.Equal(t, 1, c.Count)

	c.RegisterKill() // resets timer
	c.Update(1.9)    // still alive
	require.Equal(t, 2, c.Count)
}

func Test_Combo_XPMult_Tiers(t *testing.T) {
	c := NewCombo()
	require.Equal(t, 1.0, c.XPMult())

	for i := 0; i < 5; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 1.5, c.XPMult())

	for i := 0; i < 5; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 2.0, c.XPMult())

	for i := 0; i < 10; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 3.0, c.XPMult())
}

func Test_Combo_DamageMult_Tiers(t *testing.T) {
	c := NewCombo()
	require.Equal(t, 1.0, c.DamageMult())

	for i := 0; i < 10; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 1.10, c.DamageMult())

	for i := 0; i < 10; i++ {
		c.RegisterKill()
	}
	require.Equal(t, 1.25, c.DamageMult())
}

func Test_Combo_Label_Thresholds(t *testing.T) {
	c := NewCombo()
	require.Equal(t, "", c.Label())

	for i := 0; i < 5; i++ {
		c.RegisterKill()
	}
	require.Equal(t, "COMBO x5!", c.Label())

	for i := 0; i < 10; i++ {
		c.RegisterKill()
	}
	require.Equal(t, "MEGA KILL! x15", c.Label())

	for i := 0; i < 10; i++ {
		c.RegisterKill()
	}
	require.Equal(t, "UNSTOPPABLE! x25", c.Label())
}

func Test_Combo_Update_NoCountDoesNothing(t *testing.T) {
	c := NewCombo()
	c.Update(10.0)
	require.Equal(t, 0, c.Count)
}

func Test_Combo_Frozen_NoDecay(t *testing.T) {
	c := NewCombo()
	for i := 0; i < 10; i++ {
		c.RegisterKill()
	}
	c.Frozen = true
	c.Update(100.0) // way past timeout
	require.Equal(t, 10, c.Count, "frozen combo should not decay")
}
