package v1

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/group"

	ibcfeetypes "github.com/cosmos/ibc-go/v5/modules/apps/29-fee/types"

	"github.com/sge-network/sge/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the v1.1.0 upgrade.
const UpgradeName = "v1.1.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			group.ModuleName,
			ibcfeetypes.StoreKey,
		},
		Deleted: []string{},
	},
}
