package keeper

import (
	sdkerrors "cosmossdk.io/errors"
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

// SetPendingBet sets an pending bet
func (k Keeper) SetPendingBet(
	ctx sdk.Context,
	pendingBet *types.PendingBet,
	id uint64,
	marketUID string,
) {
	store := k.getPendingBetStore(ctx)
	b := k.cdc.MustMarshal(pendingBet)
	store.Set(types.PendingBetOfMarketKey(marketUID, id), b)
}

// IsAnyPendingBetForMarket checks if there is any pending bet for the market
func (k Keeper) IsAnyPendingBetForMarket(ctx sdk.Context, marketUID string) (thereIs bool, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingBetListOfMarketPrefix(marketUID))

	// create iterator for all existing records
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer func() {
		err = iterator.Close()
	}()

	// check if the iterator has any records
	thereIs = iterator.Valid()

	return
}

// GetPendingBets returns list of the pending bets
func (k Keeper) GetPendingBets(ctx sdk.Context) (list []types.PendingBet, err error) {
	store := k.getPendingBetStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingBet
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// RemovePendingBet removes an pending bet
func (k Keeper) RemovePendingBet(ctx sdk.Context, marketUID string, betID uint64) {
	store := k.getPendingBetStore(ctx)
	store.Delete(types.PendingBetOfMarketKey(marketUID, betID))
}

// SetSettledBet sets a settled bet
func (k Keeper) SetSettledBet(
	ctx sdk.Context,
	settledBet *types.SettledBet,
	id uint64,
	blockHeight int64,
) {
	store := k.getSettledBetStore(ctx)
	b := k.cdc.MustMarshal(settledBet)
	store.Set(types.SettledBetOfMarketKey(blockHeight, id), b)
}

// GetSettledBets returns list of the pending bets
func (k Keeper) GetSettledBets(ctx sdk.Context) (list []types.SettledBet, err error) {
	store := k.getSettledBetStore(ctx)
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

func (k Keeper) PrepareBetObject(ctx sdk.Context, creator string, props *types.WagerProps) (*types.Bet, error) {
	// Check if the value already exists
	_, isFound := k.GetBetID(ctx, props.UID)
	if isFound {
		return nil, sdkerrors.Wrapf(types.ErrDuplicateUID, "%s", props.UID)
	}

	payload := &types.WagerTicketPayload{}
	err := k.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), props.Ticket, &payload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err = payload.Validate(creator); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketValidation, "%s", err)
	}

	bet := types.NewBet(creator, props, payload.OddsType, payload.SelectedOdds)
	return bet, nil
}
