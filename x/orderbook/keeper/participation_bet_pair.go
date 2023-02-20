package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// SetParticipationBetPair sets a participation bet pair.
func (k Keeper) SetParticipationBetPair(ctx sdk.Context, bp types.ParticipationBetPair, betID uint64) {
	bpKey := types.GetParticipationBetPairKey(bp.BookID, bp.ParticipationIndex, betID)

	store := k.getParticipationBetPairStore(ctx)
	b := k.cdc.MustMarshal(&bp)
	store.Set(bpKey, b)
}
