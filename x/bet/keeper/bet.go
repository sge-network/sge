package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/bet/types"
)

// SetBet sets a specific bet in the store
func (k Keeper) SetBet(ctx sdk.Context, bet types.Bet, id uint64) {
	store := k.getBetStore(ctx)
	b := k.cdc.MustMarshal(&bet)
	store.Set(types.BetIDKey(bet.Creator, id), b)
	k.SetBetID(ctx, types.UID2ID{
		UID: bet.UID,
		ID:  id,
	})
}

// GetBet returns a bet by its UID
func (k Keeper) GetBet(ctx sdk.Context, creator string, id uint64) (val types.Bet, found bool) {
	store := k.getBetStore(ctx)

	b := store.Get(types.BetIDKey(creator, id))
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

// GetBetIDs returns list of uid2id
func (k Keeper) GetBetIDs(ctx sdk.Context) (list []types.UID2ID, err error) {
	store := k.getBetIDStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UID2ID
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// SetActiveBet sets an active bet
func (k Keeper) SetActiveBet(ctx sdk.Context, activeBet *types.ActiveBet, id uint64, sportEventUID string) {
	store := k.getActiveStore(ctx)
	b := k.cdc.MustMarshal(activeBet)
	store.Set(types.ActiveBeOfSportEventKey(sportEventUID, id), b)
}

// IsAnyActiveBetForSportevent checks if there is any active bet for the sport-event
func (k Keeper) IsAnyActiveBetForSportevent(ctx sdk.Context, sportEventUID string) (thereIs bool, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActiveBetListOfSportEventPrefix(sportEventUID))

	// create iterator for all existing records
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer func() {
		err = iterator.Close()
	}()

	// check if the iterator has any records
	thereIs = iterator.Valid()

	return
}

// GetActiveBets returns list of the active bets
func (k Keeper) GetActiveBets(ctx sdk.Context) (list []types.ActiveBet, err error) {
	store := k.getActiveStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ActiveBet
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// RemoveActiveBet removes an active bet
func (k Keeper) RemoveActiveBet(ctx sdk.Context, sportEventUID string, betID uint64) {
	store := k.getActiveStore(ctx)
	store.Delete(types.ActiveBeOfSportEventKey(sportEventUID, betID))
}

// SetSettledBet sets a settled bet
func (k Keeper) SetSettledBet(ctx sdk.Context, settledBet *types.SettledBet, id uint64, blockHeight int64) {
	store := k.getSettledStore(ctx)
	b := k.cdc.MustMarshal(settledBet)
	store.Set(types.SettledBeOfSportEventKey(blockHeight, id), b)
}

// GetSettledBets returns list of the active bets
func (k Keeper) GetSettledBets(ctx sdk.Context) (list []types.SettledBet, err error) {
	store := k.getSettledStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SettledBet
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
