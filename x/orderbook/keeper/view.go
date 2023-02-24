package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// getBookStore gets the store containing all order books.
func (k Keeper) getBookStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.BookKeyPrefix)
}

// getBookStatsStore gets the store containing the statistics.
func (k Keeper) getBookStatsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.BookStatsKeyPrefix)
}

// getParticipationStore gets the store containing all participations.
func (k Keeper) getParticipationStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.BookParticipationKeyPrefix)
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

// getHistoricalParticipationExposureStore gets the store containing all historicalparticipation exposures.
func (k Keeper) getHistoricalParticipationExposureStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.HistoricalParticipationExposureKeyPrefix)
}

// getParticipationBetPairStore gets the store containing all participation bet pair.
func (k Keeper) getParticipationBetPairStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.ParticipationBetPairKeyPrefix)
}

// getBookOddsExposureStore gets the store containing all book odds exposure.
func (k Keeper) getBookOddsExposureStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.BookOddsExposureKeyPrefix)
}

// getPayoutLockStore gets the store containing all payout locks.
func (k Keeper) getPayoutLockStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.PayoutLockKeyPrefix)
}
