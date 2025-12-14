package entity

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/stretchr/testify/require"
)

func Test_Enemy_NewNormal(t *testing.T) {
	e := NewEnemy(EnemyNormal, 10, 20)
	require.Equal(t, 10.0, e.Body.Pos.X)
	require.Equal(t, 20.0, e.Body.Pos.Y)
	require.Equal(t, EnemyNormal, e.Type)
	require.Greater(t, e.HP, 0)
	require.Equal(t, 1, e.XPDrop)
}

func Test_Enemy_NewBoss(t *testing.T) {
	e := NewEnemy(EnemyBoss, 0, 0)
	require.Equal(t, EnemyBoss, e.Type)
	require.Greater(t, e.HP, 50)
	require.Equal(t, 10, e.XPDrop)
}

func Test_Enemy_NewSwarmer(t *testing.T) {
	e := NewEnemy(EnemySwarmer, 5, 5)
	require.Equal(t, EnemySwarmer, e.Type)
	require.Equal(t, 5, e.HP)
	require.Equal(t, 3, e.Damage)
	require.Equal(t, 8.0, e.Speed)
	require.Equal(t, 1, e.XPDrop)
}

func Test_Enemy_NewTank(t *testing.T) {
	e := NewEnemy(EnemyTank, 5, 5)
	require.Equal(t, EnemyTank, e.Type)
	require.Equal(t, 60, e.HP)
	require.Equal(t, 15, e.Damage)
	require.Equal(t, 2.0, e.Speed)
	require.Equal(t, 5, e.XPDrop)
}

func Test_Enemy_NewDasher(t *testing.T) {
	e := NewEnemy(EnemyDasher, 5, 5)
	require.Equal(t, EnemyDasher, e.Type)
	require.Equal(t, 10, e.HP)
	require.Equal(t, 25, e.Damage)
	require.Equal(t, DasherIdleSpeed, e.Speed)
	require.Equal(t, 3, e.XPDrop)
	require.False(t, e.Dashing)
}

func Test_Enemy_DasherDash(t *testing.T) {
	e := NewEnemy(EnemyDasher, 0, 0)
	// Target far away — idle speed
	e.ChaseTarget(physics.Vec2{X: 100, Y: 0})
	require.Equal(t, DasherIdleSpeed, e.Speed)
	require.False(t, e.Dashing)

	// Target close — dash speed
	e.ChaseTarget(physics.Vec2{X: 10, Y: 0})
	require.Equal(t, DasherDashSpeed, e.Speed)
	require.True(t, e.Dashing)
}

func Test_Enemy_ChasePlayer(t *testing.T) {
	e := NewEnemy(EnemyNormal, 0, 0)
	playerPos := physics.Vec2{X: 10, Y: 0}
	e.ChaseTarget(playerPos)
	e.Update(1.0)
	require.Greater(t, e.Body.Pos.X, 0.0)
}

func Test_Enemy_TakeDamage(t *testing.T) {
	e := NewEnemy(EnemyNormal, 0, 0)
	initialHP := e.HP
	dead := e.TakeDamage(5)
	require.Equal(t, initialHP-5, e.HP)
	require.False(t, dead)
}

func Test_Enemy_TakeDamage_Dies(t *testing.T) {
	e := NewEnemy(EnemyNormal, 0, 0)
	dead := e.TakeDamage(9999)
	require.True(t, dead)
}

func Test_Enemy_IsAlive(t *testing.T) {
	e := NewEnemy(EnemyNormal, 0, 0)
	require.True(t, e.IsAlive())
	e.HP = 0
	require.False(t, e.IsAlive())
}

func Test_ScaleForWave_Wave1_NoChange(t *testing.T) {
	e := NewEnemy(EnemyNormal, 0, 0)
	origHP := e.HP
	origDmg := e.Damage
	origSpeed := e.Speed

	e.ScaleForWave(1)

	require.Equal(t, origHP, e.HP, "HP should not change on wave 1")
	require.Equal(t, origDmg, e.Damage, "Damage should not change on wave 1")
	require.Equal(t, origSpeed, e.Speed, "Speed should not change on wave 1")
}

func Test_ScaleForWave_Wave10(t *testing.T) {
	e := NewEnemy(EnemyNormal, 0, 0)
	// Normal: HP=15, Damage=10
	require.Equal(t, 15, e.HP)
	require.Equal(t, 10, e.Damage)

	e.ScaleForWave(10)

	// hpMult = 1 + 0.08*9 = 1.72 -> int(15 * 1.72) = int(25.8) = 25
	// dmgMult = 1 + 0.05*9 = 1.45 -> int(10 * 1.45) = int(14.5) = 14
	require.Equal(t, 25, e.HP)
	require.Equal(t, 25, e.MaxHP)
	require.Equal(t, 14, e.Damage)
}

func Test_ScaleForWave_Boss(t *testing.T) {
	e := NewEnemy(EnemyBoss, 0, 0)
	// Boss: HP=80, Damage=20
	require.Equal(t, 80, e.HP)
	require.Equal(t, 20, e.Damage)
	origSpeed := e.Speed

	e.ScaleForWave(10)

	// hpMult = 1.72 -> int(80 * 1.72) = int(137.6) = 137
	// dmgMult = 1.45 -> int(20 * 1.45) = int(29.0) = 29
	require.Equal(t, 137, e.HP)
	require.Equal(t, 137, e.MaxHP)
	require.Equal(t, 29, e.Damage)
	require.Equal(t, origSpeed, e.Speed, "Speed should not change from scaling")
}
