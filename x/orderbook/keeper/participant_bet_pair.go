package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// SetParticipantBetPair sets a participant bet pair.
func (k Keeper) SetParticipantBetPair(ctx sdk.Context, bp types.ParticipantBetPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantBetPairKeyPrefix)
	store.Set(types.GetParticipantBetPairKey(bp.BookID, bp.ParticipantNumber, bp.BetID), types.MustMarshalParticipantBetPair(k.cdc, bp))
}
