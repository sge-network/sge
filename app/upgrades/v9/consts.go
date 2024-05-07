package v9

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/sge-network/sge/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the v1.7.0 upgrade.
const UpgradeName = "v1.7.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			crisistypes.ModuleName,
			consensustypes.ModuleName,
			ibchookstypes.StoreKey,
			wasmtypes.ModuleName,
			ibcwasmtypes.ModuleName,
		},
		Deleted: []string{},
	},
}
