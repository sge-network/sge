package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// SetParticipationExposure sets a participation exposure.
func (k Keeper) SetParticipationExposure(ctx sdk.Context, pe types.ParticipationExposure) {
	peKey := types.GetParticipationExposureKey(pe.BookID, pe.OddsID, pe.ParticipationIndex)

	store := k.getParticipationExposureStore(ctx)
	b := k.cdc.MustMarshal(&pe)
	store.Set(peKey, b)

	peKey = types.GetParticipationExposureByIndexKey(pe.BookID, pe.OddsID, pe.ParticipationIndex)

	store = k.getParticipationExposureByIndexStore(ctx)
	b = k.cdc.MustMarshal(&pe)
	store.Set(peKey, b)
}

func (k Keeper) MoveToHistoricalParticipationExposure(ctx sdk.Context, pe types.ParticipationExposure) {
	peKey := types.GetHistoricalParticipationExposureKey(pe.BookID, pe.OddsID, pe.ParticipationIndex, pe.Round)

	store := k.getHistoricalParticipationExposureStore(ctx)
	b := k.cdc.MustMarshal(&pe)
	store.Set(peKey, b)

	store = k.getParticipationExposureStore(ctx)
	store.Delete(types.GetParticipationExposureKey(pe.BookID, pe.OddsID, pe.ParticipationIndex))

	store = k.getParticipationExposureByIndexStore(ctx)
	store.Delete(types.GetParticipationExposureByIndexKey(pe.BookID, pe.OddsID, pe.ParticipationIndex))
}

// GetExposureByBookAndOdds returns all exposures for a book id and odds id
func (k Keeper) GetExposureByBookAndOdds(ctx sdk.Context, bookID, oddsID string) (list []types.ParticipationExposure, err error) {
	store := k.getParticipationExposureStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetParticipationExposuresKey(bookID, oddsID))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ParticipationExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetExposureByBookAndParticipationIndex returns all exposures for a book id and participation index
func (k Keeper) GetExposureByBookAndParticipationIndex(ctx sdk.Context, bookID string, participationIndex uint64) (list []types.ParticipationExposure, err error) {
	store := k.getParticipationExposureByIndexStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetParticipationByIndexKey(bookID, participationIndex))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ParticipationExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllParticipationExposures returns all participation exposures used during genesis dump.
func (k Keeper) GetAllParticipationExposures(ctx sdk.Context) (list []types.ParticipationExposure, err error) {
	store := k.getParticipationExposureStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ParticipationExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
