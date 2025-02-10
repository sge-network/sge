package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/market/keeper"
	"github.com/sge-network/sge/x/legacy/market/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the markets
	for _, elem := range genState.MarketList {
		k.SetMarket(ctx, elem)
	}

	k.SetMarketStats(ctx, genState.Stats)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	var err error
	genesis.MarketList, err = k.GetMarkets(ctx)
	if err != nil {
		panic(err)
	}

	genesis.Stats = k.GetMarketStats(ctx)

	return genesis
}
