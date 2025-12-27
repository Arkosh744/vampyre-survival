package reward

import "math/rand"

type RewardID int

const (
	RewGlassCannon    RewardID = iota // Rare: half HP, x3 dmg
	RewDoomPact                       // Epic: -1 HP/s, +XP, +dmg
	RewBloodPrice                     // Rare: -3 HP per fire, x2 dmg
	RewSoulHarvest                    // Epic: no XP, weapon upgrade on kill
	RewBerserkerPact                  // Uncommon: +50% speed, always moving
	RewFragileEgo                     // Rare: x3 combo mults, hit resets combo
	RewGiantSlayer                    // Uncommon: x5 boss/tank, x0.5 normal
	RewVoidWalker                     // Rare: +3s invuln, freeze MaxHP
	RewGamblersFate                   // Uncommon: 50% x3 dmg, 50% miss
	RewCursedStrength                 // Rare: +20 all dmg, enemies 30% faster
	RewChainReaction                  // Rare: enemies explode on death
	RewMomentum                       // Rare: combo never expires
	RewGravityWell                    // Uncommon: pull enemies toward player
	RewMirrorImage                    // Epic: weapons fire twice
	RewWeaponForge                    // Uncommon: +10 dmg all weapons
	RewTimeWarp                       // Rare: -30% cooldown all weapons
	RewXPMagnet                       // Common: triple magnet radius
	RewVampiricAura                   // Common: +3 HP/s regen
	RewSecondWind                     // Epic: revive once at full HP
	RewFullRestore                    // Common: full heal + 25 max HP
)

type RewardRarity int

const (
	RarityCommon   RewardRarity = iota // wave 1+
	RarityUncommon                     // wave 3+
	RarityRare                         // wave 7+
	RarityEpic                         // wave 12+
)

type Reward struct {
	ID          RewardID
	Name        string
	Description string
	Rarity      RewardRarity
}

var allRewards = []Reward{
	// Risk/Reward
	{RewGlassCannon, "Glass Cannon", "Half max HP, triple all damage", RarityRare},
	{RewDoomPact, "Doom Pact", "-1 HP/s, +100% XP, +50% dmg", RarityEpic},
	{RewBloodPrice, "Blood Price", "Weapons cost 3 HP per hit, x2 damage", RarityRare},
	{RewSoulHarvest, "Soul Harvest", "No XP, weapon upgrade on kill", RarityEpic},
	{RewBerserkerPact, "Berserker Pact", "+50% speed, can't stop moving", RarityUncommon},
	{RewFragileEgo, "Fragile Ego", "x3 combo mults, hit resets combo", RarityRare},
	{RewGiantSlayer, "Giant Slayer", "x5 dmg to bosses, x0.5 to normals", RarityUncommon},
	{RewVoidWalker, "Void Walker", "+3s invuln, max HP locked", RarityRare},
	{RewGamblersFate, "Gambler's Fate", "50% for x3 dmg, 50% miss", RarityUncommon},
	{RewCursedStrength, "Cursed Strength", "+20 all weapon dmg, enemies 30% faster", RarityRare},
	// Modifiers
	{RewChainReaction, "Chain Reaction", "Enemies explode on death (25 AoE)", RarityRare},
	{RewMomentum, "Momentum", "Combo timer never expires", RarityRare},
	{RewGravityWell, "Gravity Well", "Enemies pulled toward player", RarityUncommon},
	{RewMirrorImage, "Mirror Image", "All weapons fire twice", RarityEpic},
	{RewWeaponForge, "Weapon Forge", "+10 damage to all weapons", RarityUncommon},
	{RewTimeWarp, "Time Warp", "-30% cooldown on all weapons", RarityRare},
	{RewXPMagnet, "XP Magnet", "Triple pickup magnet radius", RarityCommon},
	{RewVampiricAura, "Vampiric Aura", "+3 HP/s regen", RarityCommon},
	{RewSecondWind, "Second Wind", "Revive once at full HP on death", RarityEpic},
	{RewFullRestore, "Full Restore", "Full heal + 25 max HP", RarityCommon},
}

func AllRewards() []Reward {
	return allRewards
}

func minWaveForRarity(r RewardRarity) int {
	switch r {
	case RarityCommon:
		return 1
	case RarityUncommon:
		return 3
	case RarityRare:
		return 7
	case RarityEpic:
		return 12
	default:
		return 1
	}
}

type RewardPool struct {
	taken map[RewardID]bool
}

func NewRewardPool() *RewardPool {
	return &RewardPool{taken: make(map[RewardID]bool)}
}

// GetChoices returns 2 random rewards available for the given wave.
func (p *RewardPool) GetChoices(wave int) []Reward {
	var eligible []Reward
	for _, r := range allRewards {
		if !p.taken[r.ID] && wave >= minWaveForRarity(r.Rarity) {
			eligible = append(eligible, r)
		}
	}
	if len(eligible) == 0 {
		return nil
	}
	rand.Shuffle(len(eligible), func(i, j int) {
		eligible[i], eligible[j] = eligible[j], eligible[i]
	})
	n := 2
	if len(eligible) < n {
		n = len(eligible)
	}
	return eligible[:n]
}

// MarkTaken marks a reward as taken so it won't appear again.
func (p *RewardPool) MarkTaken(id RewardID) {
	p.taken[id] = true
}
