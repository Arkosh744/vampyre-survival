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

func Test_WaveEvent_Wave3_SwarmerRush(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 3
	ws.StartWave()
	require.Equal(t, EventSwarmerRush, ws.Event)
	require.Equal(t, "SWARMER RUSH", ws.EventName)
}

func Test_WaveEvent_Wave6_TankBrigade(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 6
	ws.StartWave()
	require.Equal(t, EventTankBrigade, ws.Event)
	require.Equal(t, "TANK BRIGADE", ws.EventName)
}

func Test_WaveEvent_Wave9_DasherBlitz(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 9
	ws.StartWave()
	require.Equal(t, EventDasherBlitz, ws.Event)
	require.Equal(t, "DASHER BLITZ", ws.EventName)
}

func Test_WaveEvent_Wave12_BloodMoon(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 12
	ws.StartWave()
	require.Equal(t, EventBloodMoon, ws.Event)
	require.Equal(t, "BLOOD MOON", ws.EventName)
}

func Test_WaveEvent_BossWave_NoEvent(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 5
	ws.StartWave()
	require.Equal(t, EventNone, ws.Event)
	require.Equal(t, "", ws.EventName)
}

func Test_CurrentSpawnInterval_Decreases(t *testing.T) {
	ws := NewWaveSpawner()

	ws.CurrentWave = 1
	interval1 := ws.CurrentSpawnInterval()

	ws.CurrentWave = 10
	interval10 := ws.CurrentSpawnInterval()

	require.Less(t, interval10, interval1, "spawn interval should decrease at higher waves")
}

func Test_CurrentSpawnInterval_MinFloor(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 200
	interval := ws.CurrentSpawnInterval()
	require.GreaterOrEqual(t, interval, MinSpawnInterval, "spawn interval should never go below MinSpawnInterval")
	require.Equal(t, MinSpawnInterval, interval)
}

func Test_WaveSpawner_Wave1_Total500(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = 1
	ws.StartWave()
	require.Equal(t, BaseEnemyCount, ws.TotalInWave)
}

func Test_SpawnBatchSize(t *testing.T) {
	require.Equal(t, 3, SpawnBatchSize)
}

func Test_WaveSpawner_FinalWave(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = FinalWave
	require.True(t, ws.IsFinalWave())
	ws.CurrentWave = FinalWave - 1
	require.False(t, ws.IsFinalWave())
}

func Test_WaveSpawner_NextWave_AtFinal_Finished(t *testing.T) {
	ws := NewWaveSpawner()
	ws.CurrentWave = FinalWave
	ws.WaveActive = true
	ws.NextWave()
	require.True(t, ws.Finished)
	require.False(t, ws.WaveActive)
	require.Equal(t, FinalWave, ws.CurrentWave, "should not increment past FinalWave")
}

func Test_WaveSpawner_NotFinished_Before15(t *testing.T) {
	ws := NewWaveSpawner()
	for i := 1; i < FinalWave; i++ {
		ws.CurrentWave = i
		ws.NextWave()
		require.False(t, ws.Finished, "should not be finished at wave %d", i)
	}
}
