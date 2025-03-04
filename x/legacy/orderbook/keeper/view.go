package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// getOrderBookStore gets the store containing all order books.
func (k Keeper) getOrderBookStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.OrderBookKeyPrefix)
}

// getOrderBookStatsStore gets the store containing the statistics.
func (k Keeper) getOrderBookStatsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.OrderBookStatsKeyPrefix)
}

// getParticipationStore gets the store containing all participations.
func (k Keeper) getParticipationStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.OrderBookParticipationKeyPrefix)
}

// getParticipationExposureStore gets the store containing all participation exposures.
func (k Keeper) getParticipationExposureStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.ParticipationExposureKeyPrefix)
}

// getParticipationExposureByIndexStore gets the store containing all participation exposures by index.
func (k Keeper) getParticipationExposureByIndexStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.ParticipationExposureByIndexKeyPrefix)
}

// getHistoricalParticipationExposureStore gets the store containing all historical participation exposures.
func (k Keeper) getHistoricalParticipationExposureStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.HistoricalParticipationExposureKeyPrefix)
}

// getParticipationBetPairStore gets the store containing all participation bet pair.
func (k Keeper) getParticipationBetPairStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.ParticipationBetPairKeyPrefix)
}

// getOrderBookOddsExposureStore gets the store containing all book odds exposure.
func (k Keeper) getOrderBookOddsExposureStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.OrderBookOddsExposureKeyPrefix)
}
