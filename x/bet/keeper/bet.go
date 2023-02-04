package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/bet/types"
)

// SetBet sets a specific bet in the store
func (k Keeper) SetBet(ctx sdk.Context, bet types.Bet, id uint64) {
	store := k.getBetStore(ctx)
	b := k.cdc.MustMarshal(&bet)
	store.Set(types.BetListByIDKey(bet.Creator, id), b)
	k.SetBetID(ctx, types.UID2ID{
		UID: bet.UID,
		ID:  id,
	})
}

// GetBet returns a bet by its UID
func (k Keeper) GetBet(ctx sdk.Context, creator string, id uint64) (val types.Bet, found bool) {
	store := k.getBetStore(ctx)

	b := store.Get(types.BetListByIDKey(creator, id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetBets returns list of bets
func (k Keeper) GetBets(ctx sdk.Context) (list []types.Bet, err error) {
	store := k.getBetStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Bet
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// SetBetID sets a specific bet id map in the store
func (k Keeper) SetBetID(ctx sdk.Context, uid2ID types.UID2ID) {
	store := k.getBetIDStore(ctx)
	b := k.cdc.MustMarshal(&uid2ID)
	store.Set(utils.StrBytes(uid2ID.UID), b)
}

// GetBetID returns a bet ID by its UID
func (k Keeper) GetBetID(ctx sdk.Context, uid string) (val types.UID2ID, found bool) {
	store := k.getBetIDStore(ctx)

	b := store.Get(utils.StrBytes(uid))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
