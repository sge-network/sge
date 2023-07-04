package subaccount

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(_ sdk.Context, _ keeper.Keeper, _ types.GenesisState) {}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(_ sdk.Context, _ keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	return genesis
}
