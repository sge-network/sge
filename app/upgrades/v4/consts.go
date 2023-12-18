package v3

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/sge-network/sge/app/upgrades"
	rewardmoduletypes "github.com/sge-network/sge/x/reward/types"
	subaccountmoduletypes "github.com/sge-network/sge/x/subaccount/types"
)

// UpgradeName defines the on-chain upgrade name for the v1.3.0 upgrade.
const UpgradeName = "v1.3.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			subaccountmoduletypes.ModuleName,
			rewardmoduletypes.ModuleName,
		},
		Deleted: []string{},
	},
}
