package world

import (
	"math"
	"math/rand"

	"github.com/arkosh/vampyre-survival/internal/entity"
)

const (
	BaseEnemyCount = 20
	WaveScaling    = 1.2
	BossEveryN     = 5
	SpawnInterval  = 0.1
)

type WaveSpawner struct {
	CurrentWave   int
	WaveActive    bool
	EnemiesToSpawn int
	TotalInWave   int
	SpawnTimer    float64
	NextIsBoss    bool
}

func NewWaveSpawner() *WaveSpawner {
	return &WaveSpawner{CurrentWave: 1}
}

func (ws *WaveSpawner) StartWave() {
	count := int(float64(BaseEnemyCount) * math.Pow(WaveScaling, float64(ws.CurrentWave-1)))
	if count < BaseEnemyCount {
		count = BaseEnemyCount
	}
	ws.TotalInWave = count
	ws.EnemiesToSpawn = count
	ws.WaveActive = true
	ws.SpawnTimer = 0
	ws.NextIsBoss = false

	if ws.IsBossWave() {
		ws.EnemiesToSpawn++
		ws.TotalInWave++
	}
}

func (ws *WaveSpawner) IsBossWave() bool {
	return ws.CurrentWave%BossEveryN == 0
}

func (ws *WaveSpawner) SpawnNext(worldCenterX, worldCenterY float64) *entity.Enemy {
	if ws.EnemiesToSpawn <= 0 {
		return nil
	}

	ws.EnemiesToSpawn--

	if ws.IsBossWave() && ws.EnemiesToSpawn == 0 {
		return entity.NewEnemy(entity.EnemyBoss, worldCenterX, worldCenterY)
	}

	typ := ws.pickEnemyType()
	return entity.NewEnemy(typ, worldCenterX, worldCenterY)
}

func (ws *WaveSpawner) pickEnemyType() entity.EnemyType {
	// Wave 1-2: normal only
	// Wave 3+: add swarmers
	// Wave 5+: add tanks
	// Wave 7+: add dashers
	types := []entity.EnemyType{entity.EnemyNormal}

	if ws.CurrentWave >= 3 {
		types = append(types, entity.EnemySwarmer, entity.EnemySwarmer) // double weight
	}
	if ws.CurrentWave >= 5 {
		types = append(types, entity.EnemyTank)
	}
	if ws.CurrentWave >= 7 {
		types = append(types, entity.EnemyDasher)
	}

	return types[rand.Intn(len(types))]
}

func (ws *WaveSpawner) AllSpawned() bool {
	return ws.EnemiesToSpawn <= 0
}

func (ws *WaveSpawner) NextWave() {
	ws.CurrentWave++
	ws.WaveActive = false
}
