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

	for _, book := range data.Books {
		keeper.SetBook(ctx, book)
	}

	for _, bp := range data.Bookparticipants {
		keeper.SetBookParticipant(ctx, bp)
	}

	for _, be := range data.BookExposures {
		keeper.SetBookOddExposure(ctx, be)
	}

	for _, pe := range data.ParticipantExposures {
		keeper.SetParticipantExposure(ctx, pe)
	}

	keeper.SetOrderBookStats(ctx, data.Stats)
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the params found in the keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	// TODO
	return &types.GenesisState{
		Params:               keeper.GetParams(ctx),
		Books:                keeper.GetAllBooks(ctx),
		Bookparticipants:     keeper.GetAllBookParticipants(ctx),
		BookExposures:        keeper.GetAllBookExposures(ctx),
		ParticipantExposures: keeper.GetAllParticipantExposures(ctx),
		Stats:                keeper.GetOrderBookStats(ctx),
	}
}

// ValidateGenesis validates the provided orderbook genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds)
func ValidateGenesis(data *types.GenesisState) error {
	// TODO

	return data.Params.Validate()
}
