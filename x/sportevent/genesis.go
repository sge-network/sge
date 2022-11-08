package sportevent

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/sportevent/keeper"
	"github.com/sge-network/sge/x/sportevent/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the sportEvent
	for _, elem := range genState.SportEventList {
		k.SetSportEvent(ctx, elem)
	}

	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	var err error
	genesis.SportEventList, err = k.GetSportEventAll(ctx)

	if err != nil {
		panic(err)
	}

	return genesis
}
