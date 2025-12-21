package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/arkosh/vampyre-survival/internal/engine"
	"github.com/arkosh/vampyre-survival/internal/entity"
	"github.com/arkosh/vampyre-survival/internal/physics"
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
	StatePaused
	StateGameOver
)

const (
	WorldWidth  = 200.0
	WorldHeight = 200.0
	TargetFPS   = 30
)

type Game struct {
	State    GameState
	Term     *ui.Terminal
	Input    *ui.Input
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
}

func NewGame() (*Game, error) {
	t, err := ui.NewTerminal()
	if err != nil {
		return nil, err
	}

	cam := engine.NewCamera(t.Width(), t.Height()-3)
	cam.SetWorldBounds(WorldWidth, WorldHeight)

	g := &Game{
		State:    StateMenu,
		Term:     t,
		Input:    ui.NewInput(),
		Camera:   cam,
		Renderer: ui.NewRenderer(t, cam),
		Menu:     ui.NewMenu(),
		Running:  true,
	}

	return g, nil
}

func (g *Game) Run() {
	defer g.Term.Close()
	defer g.Input.Close()

	ticker := time.NewTicker(time.Second / TargetFPS)
	defer ticker.Stop()

	for g.Running {
		<-ticker.C
		dt := 1.0 / float64(TargetFPS)

		g.handleInput()
		g.update(dt)
		g.render()
	}
}

func (g *Game) handleInput() {
	key := g.Input.Poll()
	if key == ui.KeyNone {
		return
	}

	switch g.State {
	case StateMenu:
		switch key {
		case ui.KeyUp:
			g.Menu.Up()
		case ui.KeyDown:
			g.Menu.Down()
		case ui.KeyEnter:
			switch g.Menu.Selected() {
			case ui.MenuStart:
				g.startGame()
			case ui.MenuQuit:
				g.Running = false
			}
		case ui.KeyQ, ui.KeyEsc:
			g.Running = false
		}

	case StateWeaponSelect:
		switch key {
		case ui.KeyUp:
			g.WeaponSelect.Up()
		case ui.KeyDown:
			g.WeaponSelect.Down()
		case ui.KeyEnter:
			g.applyWeaponSelect()
		}

	case StatePlaying:
		dir := physics.Vec2{}
		switch key {
		case ui.KeyUp, ui.KeyW:
			dir.Y = -1
		case ui.KeyDown, ui.KeyS:
			dir.Y = 1
		case ui.KeyLeft, ui.KeyA:
			dir.X = -1
		case ui.KeyRight, ui.KeyD:
			dir.X = 1
		case ui.KeyEsc:
			g.State = StatePaused
			return
		}
		if dir.Length() > 0 {
			g.Player.SetDirection(dir)
		}

	case StateLevelUp:
		switch key {
		case ui.KeyUp:
			g.LevelUp.Up()
		case ui.KeyDown:
			g.LevelUp.Down()
		case ui.KeyEnter:
			g.applyLevelUp()
		}

	case StatePaused:
		switch key {
		case ui.KeyEsc, ui.KeyEnter:
			g.State = StatePlaying
		case ui.KeyQ:
			g.State = StateMenu
		}

	case StateGameOver:
		switch key {
		case ui.KeyEnter:
			g.State = StateMenu
		case ui.KeyQ, ui.KeyEsc:
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

	// Passive regen
	if g.RegenRate > 0 && g.Player.HP < g.Player.MaxHP {
		g.RegenAccum += g.RegenRate * dt
		if g.RegenAccum >= 1.0 {
			heal := int(g.RegenAccum)
			g.RegenAccum -= float64(heal)
			g.Player.HP += heal
			if g.Player.HP > g.Player.MaxHP {
				g.Player.HP = g.Player.MaxHP
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
		p.MagnetToward(g.Player.Body.Pos, 8.0)
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
		for _, hit := range w.GetHits() {
			g.applyHit(hit)
		}
	}

	for _, e := range g.Enemies {
		if e.IsAlive() {
			e.ChaseTarget(g.Player.Body.Pos)
			e.Update(dt)

			if engine.CheckAABB(&g.Player.Body, &e.Body) {
				died := g.Player.TakeDamage(e.Damage)
				g.DamageTaken += e.Damage
				if died {
					g.State = StateGameOver
					return
				}
				g.Camera.Shake(1.5)
			}
		}
	}

	g.removeDeadEnemies()

	if !g.WaveSpawner.WaveActive {
		g.WaveSpawner.StartWave()
		// Show banner for wave events
		if g.WaveSpawner.EventName != "" {
			g.BannerText = g.WaveSpawner.EventName
			g.BannerTimer = 2.0
			g.BannerColor = ui.ColorBoldRed
		}
	}

	if g.WaveSpawner.WaveActive {
		g.WaveSpawner.SpawnTimer += dt
		if g.WaveSpawner.SpawnTimer >= g.WaveSpawner.CurrentSpawnInterval() && !g.WaveSpawner.AllSpawned() {
			g.WaveSpawner.SpawnTimer = 0
			e := g.spawnEnemyAtEdge()
			if e != nil {
				g.Enemies = append(g.Enemies, e)
			}
		}
	}

	if g.WaveSpawner.AllSpawned() && len(g.Enemies) == 0 {
		g.WaveSpawner.NextWave()
	}
}

func (g *Game) render() {
	g.Term.BeginFrame()

	switch g.State {
	case StateMenu:
		g.Menu.Draw(g.Term)

	case StateWeaponSelect:
		if g.WeaponSelect != nil {
			g.WeaponSelect.Draw(g.Term)
		}

	case StatePlaying:
		g.Renderer.DrawGround()
		g.Renderer.DrawWorldBorder(WorldWidth, WorldHeight)
		for _, p := range g.Pickups {
			g.Renderer.DrawPickup(p)
		}
		for _, e := range g.Enemies {
			g.Renderer.DrawEnemy(e)
		}
		for _, w := range g.Weapons {
			g.Renderer.DrawWeaponVisuals(w.GetVisuals())
		}
		g.Renderer.DrawEffects()
		g.Renderer.DrawPlayer(g.Player)
		g.Renderer.DrawHUD(g.Player, g.WaveSpawner.CurrentWave, g.Kills, len(g.Weapons), g.Combo.Count)

		// Combo label banner
		if label := g.Combo.Label(); label != "" {
			comboColor := ui.ColorYellow
			if g.Combo.Count >= 20 {
				comboColor = ui.ColorBoldRed
			} else if g.Combo.Count >= 10 {
				comboColor = ui.ColorMagenta
			}
			w := g.Term.Width()
			g.Term.WriteStr((w-len(label))/2, 3, label, comboColor)
		}

		// Wave event banner
		if g.BannerTimer > 0 && g.BannerText != "" {
			w := g.Term.Width()
			g.Term.WriteStr((w-len(g.BannerText))/2, 2, g.BannerText, g.BannerColor)
		}

	case StateLevelUp:
		if g.LevelUp != nil {
			g.LevelUp.Draw(g.Term)
		}

	case StatePaused:
		w := g.Term.Width()
		h := g.Term.Height()
		msg := "=== PAUSED ==="
		g.Term.WriteStr((w-len(msg))/2, h/2, msg, ui.ColorYellow)
		hint := "ESC to resume, Q to menu"
		g.Term.WriteStr((w-len(hint))/2, h/2+2, hint, ui.ColorDim)

	case StateGameOver:
		g.renderGameOver()
	}

	g.Term.EndFrame()
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
		g.Camera.Shake(0.5)

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

		xp := int(float64(e.XPDrop) * g.XPMultiplier * g.Combo.XPMult())
		if xp < 1 {
			xp = 1
		}
		leveled := g.Player.AddXP(xp)
		if leveled {
			g.showLevelUp()
		}
	}
}

func (g *Game) showLevelUp() {
	owned := g.ownedWeaponKinds()
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
		g.Player.AddMaxHP(int(upg.Value))
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
	if e.Type == entity.EnemyBoss {
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
	if e.Type == entity.EnemyBoss {
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
	w := g.Term.Width()
	h := g.Term.Height()

	title := "GAME OVER"
	g.Term.WriteStr((w-len(title))/2, h/2-6, title, ui.ColorBoldRed)

	minutes := int(g.PlayTime) / 60
	seconds := int(g.PlayTime) % 60

	stats := []struct {
		label string
		color string
	}{
		{label: fmt.Sprintf("Wave: %d", g.WaveSpawner.CurrentWave), color: ui.ColorWhite},
		{label: fmt.Sprintf("Kills: %d", g.Kills), color: ui.ColorRed},
		{label: fmt.Sprintf("Level: %d", g.Player.Level), color: ui.ColorCyan},
		{label: fmt.Sprintf("Time: %d:%02d", minutes, seconds), color: ui.ColorYellow},
		{label: fmt.Sprintf("Damage Dealt: %d", g.DamageDealt), color: ui.ColorGreen},
		{label: fmt.Sprintf("Damage Taken: %d", g.DamageTaken), color: ui.ColorMagenta},
		{label: fmt.Sprintf("Weapons: %d", len(g.Weapons)), color: ui.ColorWhite},
		{label: fmt.Sprintf("Best Streak: %d", g.Combo.BestStreak), color: ui.ColorYellow},
	}

	for i, s := range stats {
		g.Term.WriteStr((w-len(s.label))/2, h/2-3+i, s.label, s.color)
	}

	hint := "Press ENTER for menu, Q to quit"
	g.Term.WriteStr((w-len(hint))/2, h/2+5, hint, ui.ColorDim)
}
