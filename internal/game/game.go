package game

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/arkosh/vampyre-survival/internal/engine"
	"github.com/arkosh/vampyre-survival/internal/entity"
	"github.com/arkosh/vampyre-survival/internal/physics"
	"github.com/arkosh/vampyre-survival/internal/reward"
	"github.com/arkosh/vampyre-survival/internal/skill"
	"github.com/arkosh/vampyre-survival/internal/ui"
	"github.com/arkosh/vampyre-survival/internal/weapon"
	"github.com/arkosh/vampyre-survival/internal/world"
)

type GameState int

const (
	StateMenu GameState = iota
	StateWeaponSelect
	StatePlaying
	StateLevelUp
	StateWaveReward
	StatePaused
	StateGameOver
	StateVictory
)

const (
	WorldWidth  = 200.0
	WorldHeight = 200.0
	TargetFPS   = 30
	MaxWeapons  = 3
	LogicalW    = 640
	LogicalH    = 360
)

type Game struct {
	State    GameState
	escreen  *ui.EbitenScreen
	Renderer *ui.Renderer
	Camera   *engine.Camera

	Player   *entity.Player
	Enemies  []*entity.Enemy
	Weapons  []weapon.Weapon
	Pickups  []*entity.Pickup

	WaveSpawner  *world.WaveSpawner
	Upgrades     *skill.UpgradePool

	Menu         *ui.Menu
	LevelUp      *ui.LevelUpScreen
	WeaponSelect *ui.WeaponSelectScreen

	Kills       int
	DamageDealt int
	DamageTaken int
	PlayTime    float64
	Running     bool

	// Passive upgrade fields
	RegenRate       float64
	RegenAccum      float64
	LifestealAmount int
	XPMultiplier    float64
	CritChance      float64

	// Banner system
	BannerText  string
	BannerTimer float64
	BannerColor string

	Combo *Combo

	// Wave reward system
	WaveRewardScreen *ui.WaveRewardScreen
	RewardPool       *reward.RewardPool
	MagnetRadius     float64

	// Risk/Reward flags
	GlassCannonActive   bool
	DoomPactActive      bool
	BloodPriceActive    bool
	SoulHarvestActive   bool
	BerserkerPactActive bool
	FragileEgoActive    bool
	GiantSlayerActive   bool
	VoidWalkerActive    bool
	GamblerActive       bool
	CursedStrActive     bool
	EnemySpeedMult      float64

	ElderBoss *entity.Enemy

	// Modifier flags
	ChainReactionDmg  int
	MirrorImageActive bool
	GravityWellActive bool
	HasSecondWind     bool
	LastDir           physics.Vec2
}

func NewGame() *Game {
	scr := ui.NewEbitenScreen(LogicalW, LogicalH)
	cam := engine.NewCamera(LogicalW/engine.TileSize, LogicalH/engine.TileSize)
	cam.Smoothness = 0.05
	cam.SetWorldBounds(WorldWidth, WorldHeight)

	return &Game{
		State:    StateMenu,
		escreen:  scr,
		Camera:   cam,
		Renderer: ui.NewRenderer(cam),
		Menu:     ui.NewMenu(),
		Running:  true,
	}
}

// Update implements ebiten.Game interface — called every tick (~60 TPS).
func (g *Game) Update() error {
	if !g.Running {
		return ebiten.Termination
	}

	dt := 1.0 / 60.0
	if tps := ebiten.ActualTPS(); tps > 0 {
		dt = 1.0 / tps
	}

	g.handleInput()
	g.update(dt)
	return nil
}

// Draw implements ebiten.Game interface — called every frame.
func (g *Game) Draw(screen *ebiten.Image) {
	g.escreen.SetTarget(screen)
	g.render()
}

// Layout implements ebiten.Game interface — returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return LogicalW, LogicalH
}

func (g *Game) handleInput() {
	switch g.State {
	case StateMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.Menu.Up()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.Menu.Down()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch g.Menu.Selected() {
			case ui.MenuStart:
				g.startGame()
			case ui.MenuQuit:
				g.Running = false
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.Running = false
		}

	case StateWeaponSelect:
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.WeaponSelect.Up()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.WeaponSelect.Down()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.applyWeaponSelect()
		}

	case StatePlaying:
		dir := physics.Vec2{}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			dir.Y = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			dir.Y = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			dir.X = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			dir.X = 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = StatePaused
			return
		}

		if dir.Length() > 0 {
			g.LastDir = dir
			g.Player.SetDirection(dir)
		} else if g.BerserkerPactActive && g.LastDir.Length() > 0 {
			g.Player.SetDirection(g.LastDir)
		} else {
			// No movement keys held and no BerserkerPact — stop the player
			g.Player.SetDirection(physics.Vec2{})
		}

	case StateLevelUp:
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.LevelUp.Up()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.LevelUp.Down()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.applyLevelUp()
		}

	case StateWaveReward:
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.WaveRewardScreen.Up()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.WaveRewardScreen.Down()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.applyWaveReward()
		}

	case StatePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.State = StatePlaying
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.State = StateMenu
		}

	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.State = StateMenu
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.Running = false
		}

	case StateVictory:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.State = StateMenu
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.Running = false
		}
	}
}

func (g *Game) update(dt float64) {
	if g.State != StatePlaying {
		return
	}

	g.PlayTime += dt
	g.Combo.Update(dt)

	// Passive regen (with DoomPact drain)
	netRegen := g.RegenRate
	if g.DoomPactActive {
		netRegen -= 1.0
	}
	if netRegen > 0 && g.Player.HP < g.Player.MaxHP {
		g.RegenAccum += netRegen * dt
		if g.RegenAccum >= 1.0 {
			heal := int(g.RegenAccum)
			g.RegenAccum -= float64(heal)
			g.Player.HP += heal
			if g.Player.HP > g.Player.MaxHP {
				g.Player.HP = g.Player.MaxHP
			}
		}
	} else if netRegen < 0 {
		g.RegenAccum += netRegen * dt
		if g.RegenAccum <= -1.0 {
			dmg := int(-g.RegenAccum)
			g.RegenAccum += float64(dmg)
			g.Player.HP -= dmg
			if g.Player.HP <= 0 {
				g.Player.HP = 0
				if g.HasSecondWind {
					g.HasSecondWind = false
					g.Player.HP = g.Player.MaxHP
					g.Player.Invulnerable = true
					g.Player.InvulnTimer = 2.0
				} else {
					g.State = StateGameOver
					return
				}
			}
		}
	}

	// Banner decay
	if g.BannerTimer > 0 {
		g.BannerTimer -= dt
	}

	g.Renderer.UpdateEffects(dt)
	g.Player.Update(dt)
	g.clampPlayerToWorld()

	g.Camera.Follow(g.Player.Body.Pos, dt)

	// Update pickups
	for _, p := range g.Pickups {
		if p.Collected {
			continue
		}
		p.MagnetToward(g.Player.Body.Pos, g.MagnetRadius)
		p.Update(dt)
		if p.Body.Pos.DistanceTo(g.Player.Body.Pos) < 1.5 {
			p.Collect()
			if p.Type == entity.PickupHP {
				g.Player.HP += p.Value
				if g.Player.HP > g.Player.MaxHP {
					g.Player.HP = g.Player.MaxHP
				}
			}
		}
	}
	g.removeCollectedPickups()

	targets := g.buildTargets()
	for _, w := range g.Weapons {
		w.Update(dt, g.Player.Body.Pos, targets)
		hits := w.GetHits()
		if g.MirrorImageActive {
			hits = append(hits, hits...)
		}
		hasFired := len(hits) > 0
		for _, hit := range hits {
			g.applyHit(hit)
		}
		if hasFired && g.BloodPriceActive {
			g.Player.HP -= 3
			if g.Player.HP <= 0 {
				g.Player.HP = 0
				if g.HasSecondWind {
					g.HasSecondWind = false
					g.Player.HP = g.Player.MaxHP
					g.Player.Invulnerable = true
					g.Player.InvulnTimer = 2.0
				} else {
					g.State = StateGameOver
					return
				}
			}
		}
	}

	for _, e := range g.Enemies {
		if e.IsAlive() {
			e.ChaseTarget(g.Player.Body.Pos)
			if g.GravityWellActive {
				pull := g.Player.Body.Pos.Sub(e.Body.Pos).Normalize().Scale(2.0 * dt)
				e.Body.Pos = e.Body.Pos.Add(pull)
			}
			e.Update(dt)

			if engine.CheckAABB(&g.Player.Body, &e.Body) {
				died := g.Player.TakeDamage(e.Damage)
				g.DamageTaken += e.Damage
				if g.FragileEgoActive {
					g.Combo.Count = 0
				}
				if died {
					if g.HasSecondWind {
						g.HasSecondWind = false
						g.Player.HP = g.Player.MaxHP
						g.Player.Invulnerable = true
						g.Player.InvulnTimer = 2.0
					} else {
						g.State = StateGameOver
						return
					}
				}
			}
		}
	}

	// Elder Vampyre boss mechanics
	if g.ElderBoss != nil && g.ElderBoss.IsAlive() {
		if g.ElderBoss.CheckEnrage() {
			g.BannerText = "THE ELDER AWAKENS!"
			g.BannerTimer = 2.5
			g.BannerColor = ui.ColorBoldRed
		}
		if g.ElderBoss.UpdateSpawnTimer(dt) {
			count := g.ElderBoss.SpawnCount()
			for i := 0; i < count; i++ {
				minion := entity.NewEnemy(entity.EnemySwarmer, g.ElderBoss.Body.Pos.X, g.ElderBoss.Body.Pos.Y)
				minion.ScaleForWave(g.WaveSpawner.CurrentWave)
				minion.Speed *= g.EnemySpeedMult
				minion.Body.MaxSpeed *= g.EnemySpeedMult
				g.Enemies = append(g.Enemies, minion)
			}
		}
	}

	g.processChainReaction()
	g.removeDeadEnemies()

	// Check victory: final wave, all spawned, all dead
	if g.WaveSpawner.IsFinalWave() && g.WaveSpawner.AllSpawned() && len(g.Enemies) == 0 {
		g.State = StateVictory
		return
	}

	if g.WaveSpawner.IsComplete() {
		g.State = StateVictory
		return
	}

	if !g.WaveSpawner.WaveActive {
		g.WaveSpawner.StartWave()
		if g.WaveSpawner.IsFinalWave() {
			g.BannerText = "FINAL WAVE INCOMING"
			g.BannerTimer = 2.5
			g.BannerColor = ui.ColorBoldRed
		} else if g.WaveSpawner.EventName != "" {
			g.BannerText = g.WaveSpawner.EventName
			g.BannerTimer = 2.0
			g.BannerColor = ui.ColorBoldRed
		}
	}

	if g.WaveSpawner.WaveActive {
		g.WaveSpawner.SpawnTimer += dt
		if g.WaveSpawner.SpawnTimer >= g.WaveSpawner.CurrentSpawnInterval() && !g.WaveSpawner.AllSpawned() {
			g.WaveSpawner.SpawnTimer = 0
			for i := 0; i < world.SpawnBatchSize && !g.WaveSpawner.AllSpawned(); i++ {
				e := g.spawnEnemyAtEdge()
				if e != nil {
					g.Enemies = append(g.Enemies, e)
				}
			}
		}
	}

	if g.WaveSpawner.AllSpawned() && len(g.Enemies) == 0 && !g.WaveSpawner.IsFinalWave() {
		g.WaveSpawner.NextWave()
		g.showWaveReward()
	}
}

func (g *Game) render() {
	dst := g.escreen.Image()

	switch g.State {
	case StateMenu:
		g.Menu.Draw(dst)

	case StateWeaponSelect:
		if g.WeaponSelect != nil {
			g.WeaponSelect.Draw(dst)
		}

	case StatePlaying:
		g.Renderer.DrawGround(dst)
		g.Renderer.DrawWorldBorder(dst, WorldWidth, WorldHeight)
		for _, p := range g.Pickups {
			g.Renderer.DrawPickup(dst, p)
		}
		for _, e := range g.Enemies {
			g.Renderer.DrawEnemy(dst, e)
		}
		for _, w := range g.Weapons {
			g.Renderer.DrawWeaponVisuals(dst, w.GetVisuals())
		}
		g.Renderer.DrawEffects(dst)
		g.Renderer.DrawPlayer(dst, g.Player)
		var wstats []weapon.WeaponStats
		for _, w := range g.Weapons {
			wstats = append(wstats, w.Stats())
		}
		comboPct := (g.Combo.DamageMult() - 1.0) * 100
		g.Renderer.DrawHUD(dst, g.Player, g.WaveSpawner.CurrentWave, world.FinalWave, g.Kills, wstats, g.Combo.Count, comboPct)
		g.Renderer.DrawBottomBar(dst, g.WaveSpawner.CurrentWave, world.FinalWave, g.Kills, len(g.Enemies))

		if g.ElderBoss != nil && g.ElderBoss.IsAlive() {
			g.Renderer.DrawBossHP(dst, "THE ELDER VAMPYRE", g.ElderBoss.HP, g.ElderBoss.MaxHP)
		}

		// Combo label banner
		if label := g.Combo.Label(); label != "" {
			comboColor := ui.ColorYellow
			if g.Combo.Count >= 20 {
				comboColor = ui.ColorBoldRed
			} else if g.Combo.Count >= 10 {
				comboColor = ui.ColorMagenta
			}
			screenW := dst.Bounds().Dx()
			textW := len(label) * 7
			ui.DrawText(dst, (screenW-textW)/2, 65, label, ui.ColorToRGBA(comboColor))
		}

		// Wave event banner
		if g.BannerTimer > 0 && g.BannerText != "" {
			screenW := dst.Bounds().Dx()
			textW := len(g.BannerText) * 7
			ui.DrawText(dst, (screenW-textW)/2, 18, g.BannerText, ui.ColorToRGBA(g.BannerColor))
		}

	case StateLevelUp:
		if g.LevelUp != nil {
			g.LevelUp.Draw(dst)
		}

	case StateWaveReward:
		if g.WaveRewardScreen != nil {
			g.WaveRewardScreen.Draw(dst)
		}

	case StatePaused:
		screenW := dst.Bounds().Dx()
		screenH := dst.Bounds().Dy()
		msg := "=== PAUSED ==="
		textW := len(msg) * 7
		ui.DrawText(dst, (screenW-textW)/2, screenH/2, msg, ui.ColorToRGBA(ui.ColorYellow))
		hint := "ESC to resume, Q to menu"
		hintW := len(hint) * 7
		ui.DrawText(dst, (screenW-hintW)/2, screenH/2+16, hint, ui.ColorToRGBA(ui.ColorDim))

	case StateGameOver:
		g.renderGameOver()

	case StateVictory:
		g.renderVictory()
	}
}

func (g *Game) startGame() {
	g.Player = entity.NewPlayer(WorldWidth/2, WorldHeight/2)
	g.Enemies = nil
	g.Pickups = nil
	g.Weapons = nil
	g.Renderer.Effects = nil
	g.Kills = 0
	g.DamageDealt = 0
	g.DamageTaken = 0
	g.PlayTime = 0
	g.RegenRate = 0
	g.RegenAccum = 0
	g.LifestealAmount = 0
	g.XPMultiplier = 1.0
	g.CritChance = 0
	g.WaveSpawner = world.NewWaveSpawner()
	g.Combo = NewCombo()
	g.Upgrades = skill.NewUpgradePool()
	g.RewardPool = reward.NewRewardPool()
	g.MagnetRadius = 8.0
	g.EnemySpeedMult = 1.0
	g.WaveRewardScreen = nil

	// Reset all reward flags
	g.GlassCannonActive = false
	g.DoomPactActive = false
	g.BloodPriceActive = false
	g.SoulHarvestActive = false
	g.BerserkerPactActive = false
	g.FragileEgoActive = false
	g.GiantSlayerActive = false
	g.VoidWalkerActive = false
	g.GamblerActive = false
	g.CursedStrActive = false
	g.ChainReactionDmg = 0
	g.MirrorImageActive = false
	g.GravityWellActive = false
	g.HasSecondWind = false
	g.LastDir = physics.Vec2{}
	g.ElderBoss = nil

	g.WeaponSelect = ui.NewWeaponSelectScreen()
	g.State = StateWeaponSelect
}

func (g *Game) applyWeaponSelect() {
	if g.WeaponSelect == nil {
		g.State = StatePlaying
		return
	}

	choice := g.WeaponSelect.Selected()
	switch choice.Kind {
	case "sword":
		g.Weapons = []weapon.Weapon{weapon.NewSword()}
	case "projectile":
		g.Weapons = []weapon.Weapon{weapon.NewProjectileGun()}
	case "aoe":
		g.Weapons = []weapon.Weapon{weapon.NewAoE()}
	case "lightning":
		g.Weapons = []weapon.Weapon{weapon.NewLightning()}
	case "orbital":
		g.Weapons = []weapon.Weapon{weapon.NewOrbital()}
	}
	g.WeaponSelect = nil
	g.State = StatePlaying
}

func (g *Game) buildTargets() []weapon.Target {
	targets := make([]weapon.Target, 0, len(g.Enemies))
	for i, e := range g.Enemies {
		if e.IsAlive() {
			targets = append(targets, weapon.Target{Pos: e.Body.Pos, ID: i})
		}
	}
	return targets
}

func (g *Game) applyHit(hit weapon.HitResult) {
	if hit.TargetID < 0 || hit.TargetID >= len(g.Enemies) {
		return
	}
	e := g.Enemies[hit.TargetID]
	if !e.IsAlive() {
		return
	}
	dmg := int(float64(hit.Damage) * g.Combo.DamageMult())

	// Reward damage modifiers
	if g.GlassCannonActive {
		dmg *= 3
	}
	if g.DoomPactActive {
		dmg = int(float64(dmg) * 1.5)
	}
	if g.BloodPriceActive {
		dmg *= 2
	}
	if g.GiantSlayerActive {
		if e.Type == entity.EnemyBoss || e.Type == entity.EnemyTank || e.Type == entity.EnemyElderVampyre {
			dmg *= 5
		} else {
			dmg /= 2
		}
	}
	if g.GamblerActive {
		if rand.Float64() < 0.5 {
			dmg *= 3
		} else {
			dmg = 0
		}
	}

	isCrit := false
	if g.CritChance > 0 && rand.Float64() < g.CritChance {
		dmg *= 2
		isCrit = true
	}

	g.DamageDealt += dmg
	dead := e.TakeDamage(dmg)

	if isCrit {
		g.Renderer.AddDamageNumberColored(e.Body.Pos, dmg, ui.ColorBoldRed)
	} else {
		g.Renderer.AddDamageNumber(e.Body.Pos, dmg)
	}

	if dead {
		g.Kills++
		g.Player.Kills++
		g.Combo.RegisterKill()
		// Lifesteal: heal on kill
		if g.LifestealAmount > 0 {
			g.Player.HP += g.LifestealAmount
			if g.Player.HP > g.Player.MaxHP {
				g.Player.HP = g.Player.MaxHP
			}
		}

		// Death particles
		g.spawnDeathParticles(e)

		// HP drop chance
		g.trySpawnHPDrop(e)

		if g.SoulHarvestActive {
			// No XP, upgrade random weapon instead
			if len(g.Weapons) > 0 {
				g.Weapons[rand.Intn(len(g.Weapons))].Upgrade()
			}
		} else {
			xpMult := g.XPMultiplier
			if g.DoomPactActive {
				xpMult += 1.0
			}
			xp := int(float64(e.XPDrop) * xpMult * g.Combo.XPMult())
			if xp < 1 {
				xp = 1
			}
			leveled := g.Player.AddXP(xp)
			if leveled {
				g.showLevelUp()
			}
		}
	}
}

func (g *Game) showLevelUp() {
	owned := g.ownedWeaponKinds()
	if len(g.Weapons) >= MaxWeapons {
		owned = []weapon.WeaponKind{weapon.KindSword, weapon.KindProjectile, weapon.KindAoE, weapon.KindLightning, weapon.KindOrbital}
	}
	choices := g.Upgrades.GetRandomChoices(3, owned)
	if len(choices) == 0 {
		return
	}
	g.LevelUp = ui.NewLevelUpScreen(choices)
	g.State = StateLevelUp
}

func (g *Game) applyLevelUp() {
	if g.LevelUp == nil {
		g.State = StatePlaying
		return
	}

	upg := g.LevelUp.SelectedUpgrade()
	g.applyUpgrade(upg)

	g.LevelUp = nil
	g.State = StatePlaying
}

func (g *Game) applyUpgrade(upg skill.Upgrade) {
	switch upg.Type {
	case skill.UpgWeaponDamage:
		if w := g.findWeapon(upg.WeaponKind); w != nil {
			w.AddDamage(int(upg.Value))
		}
	case skill.UpgWeaponCooldown:
		if w := g.findWeapon(upg.WeaponKind); w != nil {
			w.MultiplyCooldown(upg.Value)
		}
	case skill.UpgWeaponRange:
		if w := g.findWeapon(upg.WeaponKind); w != nil {
			w.AddRange(upg.Value)
		}
	case skill.UpgNewWeapon:
		g.addNewWeapon(upg.WeaponKind)
	case skill.UpgPlayerHP:
		if !g.VoidWalkerActive {
			g.Player.AddMaxHP(int(upg.Value))
		}
	case skill.UpgPlayerSpeed:
		g.Player.AddSpeed(upg.Value)
	case skill.UpgRegen:
		g.RegenRate += upg.Value
	case skill.UpgLifesteal:
		g.LifestealAmount += int(upg.Value)
	case skill.UpgXPBonus:
		g.XPMultiplier += upg.Value
	case skill.UpgInvuln:
		g.Player.InvulnBonus += upg.Value
	case skill.UpgCrit:
		g.CritChance += upg.Value
		if g.CritChance > 0.5 {
			g.CritChance = 0.5
		}
	case skill.UpgPierce:
		for _, w := range g.Weapons {
			w.AddPierce(int(upg.Value))
		}
	}
}

func (g *Game) findWeapon(kind weapon.WeaponKind) weapon.Weapon {
	for _, w := range g.Weapons {
		if w.Kind() == kind {
			return w
		}
	}
	return nil
}

func (g *Game) addNewWeapon(kind weapon.WeaponKind) {
	if len(g.Weapons) >= MaxWeapons {
		return
	}
	for _, w := range g.Weapons {
		if w.Kind() == kind {
			return // already owned
		}
	}
	switch kind {
	case weapon.KindSword:
		g.Weapons = append(g.Weapons, weapon.NewSword())
	case weapon.KindProjectile:
		g.Weapons = append(g.Weapons, weapon.NewProjectileGun())
	case weapon.KindAoE:
		g.Weapons = append(g.Weapons, weapon.NewAoE())
	case weapon.KindLightning:
		g.Weapons = append(g.Weapons, weapon.NewLightning())
	case weapon.KindOrbital:
		g.Weapons = append(g.Weapons, weapon.NewOrbital())
	}
}

func (g *Game) ownedWeaponKinds() []weapon.WeaponKind {
	kinds := make([]weapon.WeaponKind, 0, len(g.Weapons))
	for _, w := range g.Weapons {
		kinds = append(kinds, w.Kind())
	}
	return kinds
}

func (g *Game) spawnDeathParticles(e *entity.Enemy) {
	count := 3 + rand.Intn(3) // 3-5
	if e.Type == entity.EnemyElderVampyre {
		count = 15 + rand.Intn(10) // 15-24
	} else if e.Type == entity.EnemyBoss {
		count = 8 + rand.Intn(5) // 8-12
	}
	chars := []rune{'*', 'x', '~'}
	for i := 0; i < count; i++ {
		offset := physics.Vec2{
			X: (rand.Float64() - 0.5) * 3,
			Y: (rand.Float64() - 0.5) * 3,
		}
		g.Renderer.Effects = append(g.Renderer.Effects, ui.VisualEffect{
			Pos:   e.Body.Pos.Add(offset),
			Char:  chars[rand.Intn(len(chars))],
			Color: ui.ColorRed,
			TTL:   0.3,
		})
	}
}

func (g *Game) trySpawnHPDrop(e *entity.Enemy) {
	drop := false
	if e.Type == entity.EnemyBoss || e.Type == entity.EnemyElderVampyre {
		drop = true
	} else if rand.Float64() < 0.05 {
		drop = true
	}
	if drop {
		g.Pickups = append(g.Pickups, entity.NewHPPickup(e.Body.Pos.X, e.Body.Pos.Y, 20))
	}
}

func (g *Game) removeDeadEnemies() {
	alive := g.Enemies[:0]
	for _, e := range g.Enemies {
		if e.IsAlive() {
			alive = append(alive, e)
		}
	}
	g.Enemies = alive
}

func (g *Game) removeCollectedPickups() {
	alive := g.Pickups[:0]
	for _, p := range g.Pickups {
		if !p.Collected {
			alive = append(alive, p)
		}
	}
	g.Pickups = alive
}

func (g *Game) clampPlayerToWorld() {
	if g.Player.Body.Pos.X < 1 {
		g.Player.Body.Pos.X = 1
	}
	if g.Player.Body.Pos.Y < 1 {
		g.Player.Body.Pos.Y = 1
	}
	if g.Player.Body.Pos.X > WorldWidth-2 {
		g.Player.Body.Pos.X = WorldWidth - 2
	}
	if g.Player.Body.Pos.Y > WorldHeight-2 {
		g.Player.Body.Pos.Y = WorldHeight - 2
	}
}

func (g *Game) spawnEnemyAtEdge() *entity.Enemy {
	e := g.WaveSpawner.SpawnNext(g.Player.Body.Pos.X, g.Player.Body.Pos.Y)
	if e == nil {
		return nil
	}
	e.ScaleForWave(g.WaveSpawner.CurrentWave)
	e.Speed *= g.EnemySpeedMult
	e.Body.MaxSpeed *= g.EnemySpeedMult

	// Elder Vampyre spawns at world center
	if e.Type == entity.EnemyElderVampyre {
		e.Body.Pos = physics.Vec2{X: WorldWidth / 2, Y: WorldHeight / 2}
		g.ElderBoss = e
		g.BannerText = "THE ELDER VAMPYRE"
		g.BannerTimer = 3.0
		g.BannerColor = ui.ColorBoldRed
		return e
	}

	// Blood Moon: extra HP scaling
	if g.WaveSpawner.Event == world.EventBloodMoon {
		e.HP = int(float64(e.HP) * 1.5)
		e.MaxHP = e.HP
	}

	side := rand.Intn(4)
	camW := float64(g.Camera.Width)
	camH := float64(g.Camera.Height)
	px := g.Player.Body.Pos.X
	py := g.Player.Body.Pos.Y
	offset := rand.Float64()*camW - camW/2

	switch side {
	case 0: // top
		e.Body.Pos = physics.Vec2{X: px + offset, Y: py - camH/2 - 5}
	case 1: // bottom
		e.Body.Pos = physics.Vec2{X: px + offset, Y: py + camH/2 + 5}
	case 2: // left
		e.Body.Pos = physics.Vec2{X: px - camW/2 - 5, Y: py + offset}
	case 3: // right
		e.Body.Pos = physics.Vec2{X: px + camW/2 + 5, Y: py + offset}
	}

	if e.Body.Pos.X < 1 {
		e.Body.Pos.X = 1
	}
	if e.Body.Pos.Y < 1 {
		e.Body.Pos.Y = 1
	}
	if e.Body.Pos.X > WorldWidth-2 {
		e.Body.Pos.X = WorldWidth - 2
	}
	if e.Body.Pos.Y > WorldHeight-2 {
		e.Body.Pos.Y = WorldHeight - 2
	}

	return e
}

func (g *Game) renderGameOver() {
	dst := g.escreen.Image()
	screenW := dst.Bounds().Dx()
	screenH := dst.Bounds().Dy()
	cx := screenW / 2

	// Dark background
	dst.Fill(ui.ColorToRGBA(ui.ColorDim))
	ui.DrawFilledRect(dst, float64(cx)-130, float64(screenH/2)-75, 260, 170, ui.ColorToRGBA("\033[0m"))

	title := "GAME OVER"
	titleW := len(title) * 7
	ui.DrawText(dst, (screenW-titleW)/2, screenH/2-70, title, ui.ColorToRGBA(ui.ColorBoldRed))

	minutes := int(g.PlayTime) / 60
	seconds := int(g.PlayTime) % 60

	stats := []struct {
		label string
		color string
	}{
		{label: fmt.Sprintf("Wave: %d/%d", g.WaveSpawner.CurrentWave, world.FinalWave), color: ui.ColorWhite},
		{label: fmt.Sprintf("Kills: %d", g.Kills), color: ui.ColorRed},
		{label: fmt.Sprintf("Level: %d", g.Player.Level), color: ui.ColorCyan},
		{label: fmt.Sprintf("Time: %d:%02d", minutes, seconds), color: ui.ColorYellow},
		{label: fmt.Sprintf("Damage Dealt: %d", g.DamageDealt), color: ui.ColorGreen},
		{label: fmt.Sprintf("Damage Taken: %d", g.DamageTaken), color: ui.ColorMagenta},
		{label: fmt.Sprintf("Weapons: %d", len(g.Weapons)), color: ui.ColorWhite},
		{label: fmt.Sprintf("Best Streak: %d", g.Combo.BestStreak), color: ui.ColorYellow},
	}

	for i, s := range stats {
		textW := len(s.label) * 7
		ui.DrawText(dst, (screenW-textW)/2, screenH/2-45+i*16, s.label, ui.ColorToRGBA(s.color))
	}

	hint := "Press ENTER for menu, Q to quit"
	hintW := len(hint) * 7
	ui.DrawText(dst, (screenW-hintW)/2, screenH/2+70, hint, ui.ColorToRGBA(ui.ColorDim))
}

func (g *Game) renderVictory() {
	dst := g.escreen.Image()
	screenW := dst.Bounds().Dx()
	screenH := dst.Bounds().Dy()

	// Dark background with box
	dst.Fill(ui.ColorToRGBA(ui.ColorDim))
	ui.DrawFilledRect(dst, float64(screenW/2)-140, float64(screenH/2)-95, 280, 200, ui.ColorToRGBA("\033[0m"))

	title := "=== VICTORY ==="
	titleW := len(title) * 7
	ui.DrawText(dst, (screenW-titleW)/2, screenH/2-90, title, ui.ColorToRGBA(ui.ColorGreen))

	sub := "The Elder Vampyre has been slain!"
	subW := len(sub) * 7
	ui.DrawText(dst, (screenW-subW)/2, screenH/2-74, sub, ui.ColorToRGBA(ui.ColorYellow))

	hpPct := 0.0
	if g.Player.MaxHP > 0 {
		hpPct = float64(g.Player.HP) / float64(g.Player.MaxHP)
	}
	rank := CalculateRank(RunStats{
		Time:        g.PlayTime,
		Kills:       g.Kills,
		HPPercent:   hpPct,
		Level:       g.Player.Level,
		DamageDealt: g.DamageDealt,
		BestStreak:  g.Combo.BestStreak,
	})

	rankColor := ui.ColorWhite
	switch rank {
	case RankS:
		rankColor = ui.ColorYellow
	case RankA:
		rankColor = ui.ColorGreen
	case RankB:
		rankColor = ui.ColorCyan
	case RankC:
		rankColor = ui.ColorMagenta
	case RankD:
		rankColor = ui.ColorRed
	}
	rankLine := fmt.Sprintf("RANK: %s", rank.String())
	rankW := len(rankLine) * 7
	ui.DrawText(dst, (screenW-rankW)/2, screenH/2-56, rankLine, ui.ColorToRGBA(rankColor))

	minutes := int(g.PlayTime) / 60
	seconds := int(g.PlayTime) % 60

	stats := []struct {
		label string
		color string
	}{
		{label: fmt.Sprintf("Time: %d:%02d", minutes, seconds), color: ui.ColorYellow},
		{label: fmt.Sprintf("Kills: %d", g.Kills), color: ui.ColorRed},
		{label: fmt.Sprintf("Level: %d", g.Player.Level), color: ui.ColorCyan},
		{label: fmt.Sprintf("HP: %d/%d (%.0f%%)", g.Player.HP, g.Player.MaxHP, hpPct*100), color: ui.ColorGreen},
		{label: fmt.Sprintf("Damage Dealt: %d", g.DamageDealt), color: ui.ColorWhite},
		{label: fmt.Sprintf("Best Streak: %d", g.Combo.BestStreak), color: ui.ColorMagenta},
	}

	for i, s := range stats {
		textW := len(s.label) * 7
		ui.DrawText(dst, (screenW-textW)/2, screenH/2-38+i*16, s.label, ui.ColorToRGBA(s.color))
	}

	hint := "Press ENTER for menu, Q to quit"
	hintW := len(hint) * 7
	ui.DrawText(dst, (screenW-hintW)/2, screenH/2+80, hint, ui.ColorToRGBA(ui.ColorDim))
}

func (g *Game) showWaveReward() {
	choices := g.RewardPool.GetChoices(g.WaveSpawner.CurrentWave)
	if len(choices) == 0 {
		return
	}
	g.WaveRewardScreen = ui.NewWaveRewardScreen(choices)
	g.State = StateWaveReward
}

func (g *Game) applyWaveReward() {
	if g.WaveRewardScreen == nil {
		g.State = StatePlaying
		return
	}
	rew := g.WaveRewardScreen.SelectedReward()
	g.RewardPool.MarkTaken(rew.ID)
	g.applyRewardEffect(rew)
	g.WaveRewardScreen = nil
	g.State = StatePlaying
}

func (g *Game) applyRewardEffect(rew reward.Reward) {
	switch rew.ID {
	case reward.RewGlassCannon:
		g.GlassCannonActive = true
		g.Player.MaxHP /= 2
		if g.Player.HP > g.Player.MaxHP {
			g.Player.HP = g.Player.MaxHP
		}
	case reward.RewDoomPact:
		g.DoomPactActive = true
	case reward.RewBloodPrice:
		g.BloodPriceActive = true
	case reward.RewSoulHarvest:
		g.SoulHarvestActive = true
	case reward.RewBerserkerPact:
		g.BerserkerPactActive = true
		g.Player.AddSpeed(g.Player.Body.MaxSpeed * 0.5)
	case reward.RewFragileEgo:
		g.FragileEgoActive = true
	case reward.RewGiantSlayer:
		g.GiantSlayerActive = true
	case reward.RewVoidWalker:
		g.VoidWalkerActive = true
		g.Player.InvulnBonus += 3.0
	case reward.RewGamblersFate:
		g.GamblerActive = true
	case reward.RewCursedStrength:
		g.CursedStrActive = true
		g.EnemySpeedMult = 1.3
		for _, w := range g.Weapons {
			w.AddDamage(20)
		}
	case reward.RewChainReaction:
		g.ChainReactionDmg = 25
	case reward.RewMomentum:
		g.Combo.Frozen = true
	case reward.RewGravityWell:
		g.GravityWellActive = true
	case reward.RewMirrorImage:
		g.MirrorImageActive = true
	case reward.RewWeaponForge:
		for _, w := range g.Weapons {
			w.AddDamage(10)
		}
	case reward.RewTimeWarp:
		for _, w := range g.Weapons {
			w.MultiplyCooldown(0.7)
		}
	case reward.RewXPMagnet:
		g.MagnetRadius *= 3
	case reward.RewVampiricAura:
		g.RegenRate += 3
	case reward.RewSecondWind:
		g.HasSecondWind = true
	case reward.RewFullRestore:
		g.Player.AddMaxHP(25)
		g.Player.HP = g.Player.MaxHP
	}
}

func (g *Game) processChainReaction() {
	if g.ChainReactionDmg <= 0 {
		return
	}
	for _, e := range g.Enemies {
		if e.HP != 0 {
			continue
		}
		// Dead enemy (HP==0) explodes
		for _, other := range g.Enemies {
			if !other.IsAlive() {
				continue
			}
			if e.Body.Pos.DistanceTo(other.Body.Pos) <= 3.0 {
				other.TakeDamage(g.ChainReactionDmg)
			}
		}
		e.HP = -1 // mark as processed
	}
}
