package orderbook

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
)

// InitGenesis sets the parameters for the provided keeper.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	// TODO
	keeper.SetParams(ctx, data.Params)

	for _, book := range data.BookList {
		keeper.SetBook(ctx, book)
	}

	for _, bp := range data.BookParticipationList {
		keeper.SetBookParticipation(ctx, bp)
	}

	for _, be := range data.BookExposureList {
		keeper.SetBookOddsExposure(ctx, be)
	}

	for _, pe := range data.ParticipationExposureList {
		keeper.SetParticipationExposure(ctx, pe)
	}

	keeper.SetOrderBookStats(ctx, data.Stats)
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the params found in the keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	var err error
	genesis.BookList, err = k.GetAllBooks(ctx)
	if err != nil {
		panic(err)
	}

	genesis.BookParticipationList, err = k.GetAllBookParticipations(ctx)
	if err != nil {
		panic(err)
	}

	genesis.BookExposureList, err = k.GetAllBookExposures(ctx)
	if err != nil {
		panic(err)
	}

	genesis.ParticipationExposureList, err = k.GetAllParticipationExposures(ctx)
	if err != nil {
		panic(err)
	}

	genesis.Stats = k.GetOrderBookStats(ctx)

	return genesis
}

// ValidateGenesis validates the provided orderbook genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds)
func ValidateGenesis(data *types.GenesisState) error {
	// TODO

	return data.Params.Validate()
}
