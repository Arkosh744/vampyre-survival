package skill

import (
	"math/rand"

	"github.com/arkosh/vampyre-survival/internal/weapon"
)

type UpgradeType int

const (
	UpgWeaponDamage   UpgradeType = iota
	UpgWeaponCooldown             // multiply cooldown (smaller = faster)
	UpgWeaponRange
	UpgNewWeapon
	UpgPlayerHP
	UpgPlayerSpeed
)

type Upgrade struct {
	Name        string
	Description string
	Type        UpgradeType
	WeaponKind  weapon.WeaponKind // relevant for weapon upgrades
	Value       float64
}

type UpgradePool struct {
	allWeaponKinds []weapon.WeaponKind
}

func NewUpgradePool() *UpgradePool {
	return &UpgradePool{
		allWeaponKinds: []weapon.WeaponKind{weapon.KindSword, weapon.KindProjectile, weapon.KindAoE},
	}
}

func (p *UpgradePool) GetRandomChoices(count int, ownedWeapons []weapon.WeaponKind) []Upgrade {
	pool := p.buildPool(ownedWeapons)
	if len(pool) == 0 {
		return nil
	}
	if len(pool) <= count {
		return pool
	}
	rand.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})
	return pool[:count]
}

func (p *UpgradePool) buildPool(ownedWeapons []weapon.WeaponKind) []Upgrade {
	owned := make(map[weapon.WeaponKind]bool, len(ownedWeapons))
	for _, k := range ownedWeapons {
		owned[k] = true
	}

	var pool []Upgrade

	// Weapon stat upgrades — only for owned weapons
	for _, k := range ownedWeapons {
		name := weaponDisplayName(k)
		pool = append(pool,
			Upgrade{
				Name: name + ": DMG +5", Description: "Increase " + name + " damage",
				Type: UpgWeaponDamage, WeaponKind: k, Value: 5,
			},
			Upgrade{
				Name: name + ": Speed +15%", Description: "Reduce " + name + " cooldown",
				Type: UpgWeaponCooldown, WeaponKind: k, Value: 0.85,
			},
			Upgrade{
				Name: name + ": Range +2", Description: "Increase " + name + " range",
				Type: UpgWeaponRange, WeaponKind: k, Value: 2,
			},
		)
	}

	// New weapon unlocks — only if NOT owned
	for _, k := range p.allWeaponKinds {
		if !owned[k] {
			name := weaponDisplayName(k)
			pool = append(pool, Upgrade{
				Name: "New: " + name, Description: "Unlock " + name + " weapon",
				Type: UpgNewWeapon, WeaponKind: k,
			})
		}
	}

	// Player upgrades — always available
	pool = append(pool,
		Upgrade{Name: "Max HP +20", Description: "Increase max health and heal", Type: UpgPlayerHP, Value: 20},
		Upgrade{Name: "Speed +2", Description: "Increase movement speed", Type: UpgPlayerSpeed, Value: 2},
	)

	return pool
}

func weaponDisplayName(k weapon.WeaponKind) string {
	switch k {
	case weapon.KindSword:
		return "Sword"
	case weapon.KindProjectile:
		return "Projectile"
	case weapon.KindAoE:
		return "Pulse"
	default:
		return string(k)
	}
}
