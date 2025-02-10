package orderbook

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/orderbook/keeper"
	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// InitGenesis sets the parameters for the provided keeper.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	for _, book := range data.OrderBookList {
		keeper.SetOrderBook(ctx, book)
	}

	for _, bp := range data.OrderBookParticipationList {
		keeper.SetOrderBookParticipation(ctx, bp)
	}

	for _, be := range data.OrderBookExposureList {
		keeper.SetOrderBookOddsExposure(ctx, be)
	}

	for _, pe := range data.ParticipationExposureList {
		keeper.SetParticipationExposure(ctx, pe)
	}

	for _, pe := range data.ParticipationExposureByIndexList {
		keeper.SetParticipationExposureByIndex(ctx, pe)
	}

	for _, pe := range data.HistoricalParticipationExposureList {
		keeper.SetHistoricalParticipationExposure(ctx, pe)
	}

	for _, pb := range data.ParticipationBetPairExposureList {
		betID, found := keeper.BetKeeper.GetBetID(ctx, pb.BetUID)
		if !found {
			panic("bet uid %s of the participation bet pair list not found")
		}
		keeper.SetParticipationBetPair(ctx, pb, betID.ID)
	}

	keeper.SetOrderBookStats(ctx, data.Stats)

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the params found in the keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	var err error
	genesis.OrderBookList, err = k.GetAllOrderBooks(ctx)
	if err != nil {
		panic(err)
	}

	genesis.OrderBookParticipationList, err = k.GetAllOrderBookParticipations(ctx)
	if err != nil {
		panic(err)
	}

	genesis.OrderBookExposureList, err = k.GetAllOrderBookExposures(ctx)
	if err != nil {
		panic(err)
	}

	genesis.ParticipationExposureList, err = k.GetAllParticipationExposures(ctx)
	if err != nil {
		panic(err)
	}

	genesis.ParticipationExposureByIndexList, err = k.GetAllParticipationExposures(ctx)
	if err != nil {
		panic(err)
	}

	genesis.HistoricalParticipationExposureList, err = k.GetAllHistoricalParticipationExposures(ctx)
	if err != nil {
		panic(err)
	}

	genesis.ParticipationBetPairExposureList, err = k.GetAllParticipationBetPair(ctx)
	if err != nil {
		panic(err)
	}

	genesis.Stats = k.GetOrderBookStats(ctx)

	return genesis
}
