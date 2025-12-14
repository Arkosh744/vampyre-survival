package entity

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_Player_New(t *testing.T) {
	p := NewPlayer(100, 50)
	require.Equal(t, 100.0, p.Body.Pos.X)
	require.Equal(t, 50.0, p.Body.Pos.Y)
	require.Equal(t, 100, p.HP)
	require.Equal(t, 100, p.MaxHP)
	require.Equal(t, 0, p.XP)
	require.Equal(t, 1, p.Level)
}

func Test_Player_Move(t *testing.T) {
	p := NewPlayer(0, 0)
	p.SetDirection(physics.Vec2{X: 1, Y: 0})
	p.Update(1.0)
	require.Greater(t, p.Body.Pos.X, 0.0)
}

func Test_Player_TakeDamage(t *testing.T) {
	p := NewPlayer(0, 0)
	died := p.TakeDamage(30)
	require.Equal(t, 70, p.HP)
	require.False(t, died)
	require.True(t, p.Invulnerable)
}

func Test_Player_TakeDamage_WhileInvulnerable(t *testing.T) {
	p := NewPlayer(0, 0)
	p.TakeDamage(10)
	died := p.TakeDamage(10)
	require.Equal(t, 90, p.HP)
	require.False(t, died)
}

func Test_Player_TakeDamage_Death(t *testing.T) {
	p := NewPlayer(0, 0)
	p.Invulnerable = false
	died := p.TakeDamage(100)
	require.True(t, died)
	require.Equal(t, 0, p.HP)
}

func Test_Player_AddXP_LevelUp(t *testing.T) {
	p := NewPlayer(0, 0)
	leveled := p.AddXP(100)
	require.True(t, leveled)
	require.Equal(t, 2, p.Level)
}

func Test_Player_AddXP_NoLevel(t *testing.T) {
	p := NewPlayer(0, 0)
	leveled := p.AddXP(2)
	require.False(t, leveled)
	require.Equal(t, 1, p.Level)
}

func Test_Player_XPToNextLevel(t *testing.T) {
	p := NewPlayer(0, 0)
	xp1 := p.XPToNextLevel()
	p.Level = 2
	xp2 := p.XPToNextLevel()
	require.Greater(t, xp2, xp1)
}

func Test_Player_InvulnerabilityExpires(t *testing.T) {
	p := NewPlayer(0, 0)
	p.TakeDamage(10)
	require.True(t, p.Invulnerable)
	p.InvulnTimer = 0
	p.Update(0.1)
	require.False(t, p.Invulnerable)
}

func Test_Player_AddMaxHP(t *testing.T) {
	p := NewPlayer(0, 0)
	p.AddMaxHP(20)
	require.Equal(t, 120, p.MaxHP)
	require.Equal(t, 120, p.HP) // heals on upgrade
}

func Test_Player_AddMaxHP_WhenDamaged(t *testing.T) {
	p := NewPlayer(0, 0)
	p.HP = 50
	p.AddMaxHP(20)
	require.Equal(t, 120, p.MaxHP)
	require.Equal(t, 70, p.HP) // 50 + 20
}

func Test_Player_AddSpeed(t *testing.T) {
	p := NewPlayer(0, 0)
	oldSpeed := p.Body.MaxSpeed
	p.AddSpeed(2.0)
	require.Equal(t, oldSpeed+2.0, p.Body.MaxSpeed)
}

func Test_Player_InvulnBonus(t *testing.T) {
	// Player without bonus
	p1 := NewPlayer(0, 0)
	p1.TakeDamage(10)
	require.True(t, p1.Invulnerable)
	baseTimer := p1.InvulnTimer
	require.Equal(t, InvulnDuration, baseTimer)

	// Player with InvulnBonus=0.5
	p2 := NewPlayer(0, 0)
	p2.InvulnBonus = 0.5
	p2.TakeDamage(10)
	require.True(t, p2.Invulnerable)
	bonusTimer := p2.InvulnTimer
	require.Equal(t, InvulnDuration+0.5, bonusTimer)

	// Bonus timer should be longer
	require.Greater(t, bonusTimer, baseTimer, "InvulnBonus should increase invulnerability duration")
}
