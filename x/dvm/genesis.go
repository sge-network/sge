package dvm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/keeper"
	"github.com/sge-network/sge/x/dvm/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set the Key from Genesis

	if genState.PublicKeys != nil {
		k.SetPublicKeys(ctx, *genState.PublicKeys)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	// load the default Params
	genesis.Params = k.GetParams(ctx)

	// load the public keys
	keys, found := k.GetPublicKeys(ctx)
	if found {
		genesis.PublicKeys = &keys
	}

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
