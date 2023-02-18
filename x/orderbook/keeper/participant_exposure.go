package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetExposureByBookAndOdds returns all exposures for a book id and odds id
func (k Keeper) GetExposureByBookAndOdds(ctx sdk.Context, bookID, oddsID string) (pes []types.ParticipantExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.GetParticipantExposuresKey(bookID, oddsID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pe := types.MustUnmarshalParticipantExposure(k.cdc, iterator.Value())
		pes = append(pes, pe)
	}

	return pes
}

// GetExposureByBookAndParticipantNumber returns all exposures for a book id and participant number
func (k Keeper) GetExposureByBookAndParticipantNumber(ctx sdk.Context, bookID string, pn uint64) (pes []types.ParticipantExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureByPNKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.GetParticipantByPNKey(bookID, pn))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pe := types.MustUnmarshalParticipantExposure(k.cdc, iterator.Value())
		pes = append(pes, pe)
	}

	return pes
}

// IterateAllParticipantExposures iterates through all of the participant exposures.
func (k Keeper) IterateAllParticipantExposures(ctx sdk.Context, cb func(pe types.ParticipantExposure) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.ParticipantExposureKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pe := types.MustUnmarshalParticipantExposure(k.cdc, iterator.Value())
		if cb(pe) {
			break
		}
	}
}

// GetAllParticipantExposures returns all participant exposures used during genesis dump.
func (k Keeper) GetAllParticipantExposures(ctx sdk.Context) (pes []types.ParticipantExposure) {
	k.IterateAllParticipantExposures(ctx, func(pe types.ParticipantExposure) bool {
		pes = append(pes, pe)
		return false
	})

	return pes
}

// SetParticipantExposure sets a participant exposure.
func (k Keeper) SetParticipantExposure(ctx sdk.Context, pe types.ParticipantExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureKeyPrefix)
	store.Set(types.GetParticipantExposureKey(pe.BookID, pe.OddsID, pe.ParticipantNumber), types.MustMarshalParticipantExposure(k.cdc, pe))

	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureByPNKeyPrefix)
	store.Set(types.GetParticipantExposureByPNKey(pe.BookID, pe.OddsID, pe.ParticipantNumber), types.MustMarshalParticipantExposure(k.cdc, pe))
}

func (k Keeper) MoveToHistoricalParticipantExposure(ctx sdk.Context, pe types.ParticipantExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HistoricalParticipantExposureKeyPrefix)
	store.Set(types.GetHistoricalParticipantExposureKey(pe.BookID, pe.OddsID, pe.ParticipantNumber, pe.Round), types.MustMarshalParticipantExposure(k.cdc, pe))

	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureKeyPrefix)
	store.Delete(types.GetParticipantExposureKey(pe.BookID, pe.OddsID, pe.ParticipantNumber))

	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureByPNKeyPrefix)
	store.Delete(types.GetParticipantExposureByPNKey(pe.BookID, pe.OddsID, pe.ParticipantNumber))
}
