package world

import (
	"math"
	"math/rand"

	"github.com/arkosh/vampyre-survival/internal/entity"
)

const (
	BaseEnemyCount    = 500
	WaveScaling       = 1.2
	BossEveryN        = 5
	BaseSpawnInterval = 0.05
	MinSpawnInterval  = 0.008
	SpawnBatchSize    = 3
)

type WaveEvent int

const (
	EventNone WaveEvent = iota
	EventSwarmerRush
	EventTankBrigade
	EventDasherBlitz
	EventBloodMoon
)

type WaveSpawner struct {
	CurrentWave   int
	WaveActive    bool
	EnemiesToSpawn int
	TotalInWave   int
	SpawnTimer    float64
	NextIsBoss    bool
	Event         WaveEvent
	EventName     string
}

func NewWaveSpawner() *WaveSpawner {
	return &WaveSpawner{CurrentWave: 1}
}

func (ws *WaveSpawner) StartWave() {
	count := int(float64(BaseEnemyCount) * math.Pow(WaveScaling, float64(ws.CurrentWave-1)))
	if count < BaseEnemyCount {
		count = BaseEnemyCount
	}

	ws.Event, ws.EventName = ws.determineEvent()

	// Apply event count modifier
	switch ws.Event {
	case EventSwarmerRush:
		count = int(float64(count) * 1.5)
	case EventTankBrigade:
		count = int(float64(count) * 0.6)
	case EventBloodMoon:
		count *= 2
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

func (ws *WaveSpawner) determineEvent() (WaveEvent, string) {
	if ws.IsBossWave() || ws.CurrentWave%3 != 0 {
		return EventNone, ""
	}
	switch {
	case ws.CurrentWave%12 == 0:
		return EventBloodMoon, "BLOOD MOON"
	case ws.CurrentWave%9 == 0:
		return EventDasherBlitz, "DASHER BLITZ"
	case ws.CurrentWave%6 == 0:
		return EventTankBrigade, "TANK BRIGADE"
	default:
		return EventSwarmerRush, "SWARMER RUSH"
	}
}

func (ws *WaveSpawner) CurrentSpawnInterval() float64 {
	interval := BaseSpawnInterval / (1.0 + 0.05*float64(ws.CurrentWave-1))
	if interval < MinSpawnInterval {
		interval = MinSpawnInterval
	}
	return interval
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
	// Event overrides
	switch ws.Event {
	case EventSwarmerRush:
		return entity.EnemySwarmer
	case EventTankBrigade:
		if rand.Float64() < 0.7 {
			return entity.EnemyTank
		}
		return entity.EnemyNormal
	case EventDasherBlitz:
		if rand.Float64() < 0.6 {
			return entity.EnemyDasher
		}
		return entity.EnemySwarmer
	case EventBloodMoon:
		all := []entity.EnemyType{entity.EnemyNormal, entity.EnemySwarmer, entity.EnemyTank, entity.EnemyDasher}
		return all[rand.Intn(len(all))]
	}

	// Default wave-based mixing
	types := []entity.EnemyType{entity.EnemyNormal}
	if ws.CurrentWave >= 3 {
		types = append(types, entity.EnemySwarmer, entity.EnemySwarmer)
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
