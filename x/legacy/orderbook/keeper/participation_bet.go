package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// SetParticipationBetPair sets a participation bet pair.
func (k Keeper) SetParticipationBetPair(ctx sdk.Context, bp types.ParticipationBetPair, betID uint64) {
	bpKey := types.GetParticipationBetPairKey(bp.OrderBookUID, bp.ParticipationIndex, betID)

	store := k.getParticipationBetPairStore(ctx)
	b := k.cdc.MustMarshal(&bp)
	store.Set(bpKey, b)
}

// GetAllParticipationBetPair returns all participation bet pairs used during genesis dump.
func (k Keeper) GetAllParticipationBetPair(
	ctx sdk.Context,
) (list []types.ParticipationBetPair, err error) {
	store := k.getParticipationBetPairStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ParticipationBetPair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
