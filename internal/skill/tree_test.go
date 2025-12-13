package skill

import (
	"testing"

	"github.com/arkosh/vampyre-survival/internal/weapon"
	"github.com/stretchr/testify/require"
)

func Test_UpgradePool_New(t *testing.T) {
	p := NewUpgradePool()
	require.NotNil(t, p)
}

func Test_UpgradePool_WithOneWeapon(t *testing.T) {
	p := NewUpgradePool()
	choices := p.GetRandomChoices(3, []weapon.WeaponKind{weapon.KindSword})
	require.Equal(t, 3, len(choices))
}

func Test_UpgradePool_OnlyOwnedWeaponStats(t *testing.T) {
	p := NewUpgradePool()
	owned := []weapon.WeaponKind{weapon.KindSword}
	choices := p.GetRandomChoices(20, owned) // get all

	for _, c := range choices {
		if c.Type == UpgWeaponDamage || c.Type == UpgWeaponCooldown || c.Type == UpgWeaponRange {
			require.Equal(t, weapon.KindSword, c.WeaponKind)
		}
	}
}

func Test_UpgradePool_NewWeaponOnlyIfNotOwned(t *testing.T) {
	p := NewUpgradePool()
	owned := []weapon.WeaponKind{weapon.KindSword, weapon.KindProjectile, weapon.KindAoE}
	choices := p.GetRandomChoices(20, owned)

	for _, c := range choices {
		require.NotEqual(t, UpgNewWeapon, c.Type, "should not offer new weapon when all owned")
	}
}

func Test_UpgradePool_OffersNewWeapon(t *testing.T) {
	p := NewUpgradePool()
	owned := []weapon.WeaponKind{weapon.KindSword}
	choices := p.GetRandomChoices(20, owned)

	hasNew := false
	for _, c := range choices {
		if c.Type == UpgNewWeapon {
			hasNew = true
			require.NotEqual(t, weapon.KindSword, c.WeaponKind)
		}
	}
	require.True(t, hasNew)
}

func Test_UpgradePool_AlwaysHasPlayerUpgrades(t *testing.T) {
	p := NewUpgradePool()
	choices := p.GetRandomChoices(20, []weapon.WeaponKind{weapon.KindSword})

	hasHP := false
	hasSpeed := false
	for _, c := range choices {
		if c.Type == UpgPlayerHP {
			hasHP = true
		}
		if c.Type == UpgPlayerSpeed {
			hasSpeed = true
		}
	}
	require.True(t, hasHP)
	require.True(t, hasSpeed)
}
