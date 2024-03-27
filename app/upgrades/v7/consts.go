package v7

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/sge-network/sge/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the v1.5.3 upgrade.
const UpgradeName = "v1.5.3"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
