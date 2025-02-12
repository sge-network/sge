package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// SetParticipationExposure sets a participation exposure.
func (k Keeper) SetParticipationExposure(ctx sdk.Context, pe types.ParticipationExposure) {
	peKey := types.GetParticipationExposureKey(pe.OrderBookUID, pe.OddsUID, pe.ParticipationIndex)

	store := k.getParticipationExposureStore(ctx)
	b := k.cdc.MustMarshal(&pe)
	store.Set(peKey, b)

	k.SetParticipationExposureByIndex(ctx, pe)
}

// removeParticipationExposure removes a participation exposure.
func (k Keeper) removeParticipationExposure(ctx sdk.Context, pe types.ParticipationExposure) {
	store := k.getParticipationExposureStore(ctx)
	store.Delete(types.GetParticipationExposureKey(pe.OrderBookUID, pe.OddsUID, pe.ParticipationIndex))
}

// GetExposureByOrderBookAndOdds returns all exposures for an order book uid and odds uid.
func (k Keeper) GetExposureByOrderBookAndOdds(
	ctx sdk.Context,
	bookUID, oddsUID string,
) (list []types.ParticipationExposure, err error) {
	store := k.getParticipationExposureStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetParticipationExposuresKey(bookUID, oddsUID))

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

// GetExposureByOrderBook returns all exposures for an order book uid.
func (k Keeper) GetExposureByOrderBook(
	ctx sdk.Context,
	bookUID string,
) (peMap map[uint64]map[string]*types.ParticipationExposure, err error) {
	peMap = make(map[uint64]map[string]*types.ParticipationExposure)
	store := k.getParticipationExposureStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetParticipationExposuresByOrderBookKey(bookUID))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ParticipationExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if _, ok := peMap[val.ParticipationIndex]; !ok {
			peMap[val.ParticipationIndex] = map[string]*types.ParticipationExposure{val.OddsUID: &val}
		} else {
			peMap[val.ParticipationIndex][val.OddsUID] = &val
		}
	}

	return
}

// GetAllParticipationExposures returns all participation exposures used during genesis dump.
func (k Keeper) GetAllParticipationExposures(
	ctx sdk.Context,
) (list []types.ParticipationExposure, err error) {
	store := k.getParticipationExposureStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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

// SetParticipationExposureByIndex sets a participation exposure by index.
func (k Keeper) SetParticipationExposureByIndex(ctx sdk.Context, pe types.ParticipationExposure) {
	peKey := types.GetParticipationExposureByIndexKey(
		pe.OrderBookUID,
		pe.OddsUID,
		pe.ParticipationIndex,
	)

	store := k.getParticipationExposureByIndexStore(ctx)
	b := k.cdc.MustMarshal(&pe)
	store.Set(peKey, b)
}

// GetExposureByOrderBookAndParticipationIndex returns all exposures for an order book uid and participation index.
func (k Keeper) GetExposureByOrderBookAndParticipationIndex(
	ctx sdk.Context,
	bookUID string,
	participationIndex uint64,
) (list []types.ParticipationExposure, err error) {
	store := k.getParticipationExposureByIndexStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(
		store,
		types.GetParticipationByIndexKey(bookUID, participationIndex),
	)

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

// removeParticipationExposureByIndex removes a participation exposure by index.
func (k Keeper) removeParticipationExposureByIndex(ctx sdk.Context, pe types.ParticipationExposure) {
	store := k.getParticipationExposureByIndexStore(ctx)
	store.Delete(
		types.GetParticipationExposureByIndexKey(pe.OrderBookUID, pe.OddsUID, pe.ParticipationIndex),
	)
}

// SetHistoricalParticipationExposure sets a historical participation exposure.
func (k Keeper) SetHistoricalParticipationExposure(ctx sdk.Context, pe types.ParticipationExposure) {
	peKey := types.GetHistoricalParticipationExposureKey(
		pe.OrderBookUID,
		pe.OddsUID,
		pe.ParticipationIndex,
		pe.Round,
	)

	store := k.getHistoricalParticipationExposureStore(ctx)
	b := k.cdc.MustMarshal(&pe)
	store.Set(peKey, b)
}

// GetAllHistoricalParticipationExposures returns all participation exposures used during genesis dump.
func (k Keeper) GetAllHistoricalParticipationExposures(
	ctx sdk.Context,
) (list []types.ParticipationExposure, err error) {
	store := k.getHistoricalParticipationExposureStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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

// MoveToHistoricalParticipationExposure removes the participation exposures and indices
// and sets historical participation exposures.
func (k Keeper) MoveToHistoricalParticipationExposure(ctx sdk.Context, pe types.ParticipationExposure) {
	k.SetHistoricalParticipationExposure(ctx, pe)
	k.removeParticipationExposure(ctx, pe)
	k.removeParticipationExposureByIndex(ctx, pe)
}
