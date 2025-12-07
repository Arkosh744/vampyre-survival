package world

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/entity"
	"github.com/stretchr/testify/require"
)

func Test_WaveSpawner_New(t *testing.T) {
	ws := NewWaveSpawner()
	require.Equal(t, 1, ws.CurrentWave)
	require.Equal(t, false, ws.WaveActive)
}

func Test_WaveSpawner_StartWave(t *testing.T) {
	ws := NewWaveSpawner()
	ws.StartWave()
	require.True(t, ws.WaveActive)
	require.Greater(t, ws.EnemiesToSpawn, 0)
}

func Test_WaveSpawner_WaveScaling(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 1
	ws.StartWave()
	count1 := ws.EnemiesToSpawn

	ws2 := NewWaveSpawner()
	ws2.CurrentWave = 5
	ws2.StartWave()
	count5 := ws2.EnemiesToSpawn

	require.Greater(t, count5, count1)
}

func Test_WaveSpawner_BossWave(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 5
	require.True(t, ws.IsBossWave())
	ws.CurrentWave = 3
	require.False(t, ws.IsBossWave())
}

func Test_WaveSpawner_SpawnEnemy(t *testing.T) {
	ws := NewWaveSpawner()
	ws.StartWave()
	e := ws.SpawnNext(100, 100)
	require.NotNil(t, e)
	require.Equal(t, entity.EnemyNormal, e.Type)
	require.Equal(t, ws.EnemiesToSpawn, ws.TotalInWave-1)
}

func Test_WaveSpawner_BossWaveSpawnsBoss(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 5
	ws.StartWave()

	var lastEnemy *entity.Enemy
	for ws.EnemiesToSpawn > 0 {
		lastEnemy = ws.SpawnNext(50, 50)
	}
	require.NotNil(t, lastEnemy)
	require.Equal(t, entity.EnemyBoss, lastEnemy.Type)
}

func Test_WaveSpawner_WaveComplete(t *testing.T) {
	ws := NewWaveSpawner()
	ws.StartWave()
	for ws.EnemiesToSpawn > 0 {
		ws.SpawnNext(0, 0)
	}
	require.True(t, ws.AllSpawned())
}
