<div align="center">

```
██╗   ██╗ █████╗ ███╗   ███╗██████╗ ██╗   ██╗██████╗ ███████╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗╚██╗ ██╔╝██╔══██╗██╔════╝
██║   ██║███████║██╔████╔██║██████╔╝ ╚████╔╝ ██████╔╝█████╗
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝   ╚██╔╝  ██╔══██╗██╔══╝
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║        ██║   ██║  ██║███████╗
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝        ╚═╝   ╚═╝  ╚═╝╚══════╝
    ███████╗██╗   ██╗██████╗ ██╗   ██╗██╗██╗   ██╗ █████╗ ██╗
    ██╔════╝██║   ██║██╔══██╗██║   ██║██║██║   ██║██╔══██╗██║
    ███████╗██║   ██║██████╔╝██║   ██║██║██║   ██║███████║██║
    ╚════██║██║   ██║██╔══██╗╚██╗ ██╔╝██║╚██╗ ██╔╝██╔══██║██║
    ███████║╚██████╔╝██║  ██║ ╚████╔╝ ██║ ╚████╔╝ ██║  ██║███████╗
    ╚══════╝ ╚═════╝ ╚═╝  ╚═╝  ╚═══╝  ╚═╝  ╚═══╝  ╚═╝  ╚═╝╚══════╝
```

### 🧛 Survive the night. Slay the horde

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Ebitengine](https://img.shields.io/badge/Ebitengine-v2.9-4B275F?style=for-the-badge)](https://ebitengine.org)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-0078D4?style=for-the-badge)](#-quick-start)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

</div>

---

## 🎮 What is this?

**Vampyre Survival** is a pixel-art arena survival roguelike inspired by *Vampire Survivors*. Fight through **15 waves** of increasingly deadly enemies, collect weapons and upgrades, manage risk-reward choices, and face the **Elder Vampyre** in an epic final showdown. Built with Go and [Ebitengine](https://ebitengine.org).

---

## ⚔️ Arsenal

Five auto-attacking weapons — unlock up to 3 at once and upgrade them as you level up:

| | Weapon | Style | What it does |
|---|--------|-------|-------------|
| 🗡️ | **Sword** | Melee | Slashes all enemies in range. Fast, reliable, deadly up close |
| 🔵 | **Bolt** | Projectile | Fires piercing bolts at the nearest enemy. Sniper of the bunch |
| 💜 | **Pulse** | AoE | Emits shockwaves around you. Perfect crowd control |
| ⚡ | **Zap** | Chain Lightning | Strikes one enemy, then chains to up to 7 more with decaying damage |
| 🔮 | **Orb** | Orbital | Up to 5 orbs orbit around you, shredding anything they touch |

---

## 👹 Enemies

Six enemy types with stats that scale every wave (+12% HP, +8% DMG, +7.5% SPD):

| | Type | Threat |
|---|------|--------|
| 💀 | **Normal** | Balanced all-rounder. The bread and butter of the horde |
| 🐀 | **Swarmer** | Fragile but blazing fast. They come in packs |
| 🛡️ | **Tank** | Slow-moving wall of HP. Hits like a truck |
| 💨 | **Dasher** | Charges at you with a sudden burst of speed (40.0!) |
| 👿 | **Boss** | Appears every 5 waves. Tougher, meaner, worth more XP |
| 🧛 | **Elder Vampyre** | The final boss. 800 HP, spawns minions, enrages at 50% HP |

### 🧛 Elder Vampyre — Final Boss

> *Phase 1:* Spawns 3 minions every 3 seconds while hunting you down.
>
> *Phase 2 (≤50% HP):* **ENRAGE** — speed doubles, damage jumps to 45, spawns 5 minions every 1.5 seconds. Good luck.

---

## 🌊 Wave System

15 waves stand between you and victory. Every wave brings more enemies, faster spawns, and tougher stats.

| Wave | Event | What happens |
|------|-------|-------------|
| 3, 9 | 🐀 **Swarmer Rush** | Only swarmers, +50% enemy count |
| 5, 10 | 👿 **Boss Wave** | A boss joins the fight |
| 6 | 🛡️ **Tank Brigade** | 70% tanks. Bring your big guns |
| 9 | 💨 **Dasher Blitz** | 60% dashers + 40% swarmers. Pure chaos |
| 12 | 🌑 **Blood Moon** | All enemy types, x2 enemy count |
| 15 | 🧛 **Elder Vampyre** | The final showdown |

---

## 🎁 Wave Rewards

After each wave, choose from **20 unique rewards** across 4 rarities. Some are pure power — others demand a price.

### ⚠️ Risk / Reward

| Reward | Rarity | Upside | Downside |
|--------|--------|--------|----------|
| 🔪 Glass Cannon | Rare | x3 damage | Half your HP |
| 💀 Doom Pact | Epic | +100% XP, +50% damage | Lose 1 HP/sec |
| 🩸 Blood Price | Rare | x2 weapon damage | Weapons cost 3 HP per shot |
| 👻 Soul Harvest | Epic | Weapons upgrade on every kill | No XP gain |
| 🏃 Berserker Pact | Uncommon | +50% speed | Can't stop moving |
| 🎭 Fragile Ego | Rare | x3 combo multipliers | Getting hit resets combo |
| 🗡️ Giant Slayer | Uncommon | x5 damage to bosses/tanks | x0.5 to normals |
| 🌀 Void Walker | Rare | +3s invulnerability | Max HP frozen |
| 🎲 Gambler's Fate | Uncommon | 50% chance for x3 damage | 50% chance to miss |
| 💪 Cursed Strength | Rare | +20 damage to all weapons | Enemies 30% faster |

### 💎 Power-Ups

| Reward | Rarity | Effect |
|--------|--------|--------|
| 💥 Chain Reaction | Rare | Enemies explode on death |
| ⏳ Momentum | Rare | Combo never expires |
| 🌀 Gravity Well | Uncommon | Enemies are pulled toward you |
| 🪞 Mirror Image | Epic | All weapons fire twice |
| 🔨 Weapon Forge | Uncommon | +10 damage to all weapons |
| ⏰ Time Warp | Rare | -30% cooldown on all weapons |
| 🧲 XP Magnet | Common | Triple pickup radius |
| 🧛 Vampiric Aura | Common | +3 HP/sec regeneration |
| 💚 Second Wind | Epic | Revive once with full HP |
| ❤️ Full Restore | Common | Full heal + 25 max HP |

---

## 🔥 Combo System

Chain kills to build your combo multiplier:

| Combo | Title | XP Bonus |
|-------|-------|----------|
| 5+ | **COMBO x5!** | x1.5 |
| 10+ | **MEGA KILL!** | x2.0 |
| 20+ | **UNSTOPPABLE!** | x3.0 |

Every kill also increases your damage by +0.5%. Stop killing for 2 seconds and the combo resets. Keep the pressure on!

---

## 📈 Level-Up Upgrades

Earn XP, level up, choose from **12 upgrades**:

> 🗡️ Weapon Damage ・ ⏱️ Weapon Cooldown ・ 📏 Weapon Range ・ 🆕 New Weapon
>
> ❤️ Max HP ・ 🏃 Speed ・ 💚 Regen ・ 🧛 Lifesteal
>
> ✨ XP Bonus ・ 🛡️ Invulnerability ・ 💥 Crit Chance ・ 🔱 Pierce

---

## 🏆 Rank System

After defeating the Elder Vampyre, you receive a rank based on your performance:

| Rank | Score | You are... |
|------|-------|-----------|
| **S** | 85+ | A legend among hunters |
| **A** | 65+ | A seasoned slayer |
| **B** | 45+ | A worthy fighter |
| **C** | 25+ | A survivor |
| **D** | <25 | Barely alive |

Scoring: **Time** (40 pts) + **HP remaining** (30 pts) + **Kill count** (30 pts)

---

## 📸 Screenshots

> *Coming soon — pixel-art gameplay screenshots and GIFs*

---

## 🚀 Quick Start

```bash
# Clone & build
git clone https://github.com/your-username/vampyre-survival.git
cd vampyre-survival

# Option 1: Build with code signing (Windows)
make gen-certs && make build

# Option 2: Build for Linux
make build-linux

# Option 3: Cross-compile for Windows (no signing)
make build-windows
```

Then just run the binary and survive!

### 🎮 Controls

| Key | Action |
|-----|--------|
| **W A S D** | Move |
| **Auto** | Weapons fire automatically at the nearest enemy |

---

## 🔧 Requirements

- **Go 1.24+**
- **Linux native build:** `libx11-dev libgl1-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libxxf86vm-dev`
- **Code signing (optional):** `osslsigncode`

```bash
make build          # Build + sign (Windows)
make build-linux    # Native Linux build
make build-windows  # Cross-compile for Windows
make test           # Run all tests (~155 tests)
make clean          # Remove binaries
```

---

<div align="center">

**Built with** 🩸 **and** [Go](https://go.dev) **+** [Ebitengine](https://ebitengine.org)

*© 2025 arktrix games*

</div>
