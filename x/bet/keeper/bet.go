package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/bet/types"
)

// SetBet sets a specific bet in the store
func (k Keeper) SetBet(ctx sdk.Context, bet types.Bet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetListPrefix)
	b := k.cdc.MustMarshal(&bet)
	store.Set(utils.StrBytes(bet.UID), b)

}

// GetBet returns a bet by its UID
func (k Keeper) GetBet(ctx sdk.Context, uid string) (val types.Bet, found bool) {
	store := k.getBetStore(ctx)

	b := store.Get(utils.StrBytes(uid))
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
