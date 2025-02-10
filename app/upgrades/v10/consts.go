package v10

import (
	store "cosmossdk.io/store/types"
	circuittypes "cosmossdk.io/x/circuit/types"

	"github.com/sge-network/sge/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the v1.8.0 upgrade.
const UpgradeName = "v1.8.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			circuittypes.ModuleName,
		},
		Deleted: []string{},
	},
}
