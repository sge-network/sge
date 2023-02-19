package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// SetDeposit sets a deposit.
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	depoistKey := types.GetDepositKey(deposit.Creator, deposit.SportEventUID, deposit.ParticipantID)

	store := k.getDepositsStore(ctx)
	b := k.cdc.MustMarshal(&deposit)
	store.Set(depoistKey, b)
}

// GetDeposit returns a specific deposit.
func (k Keeper) GetDeposit(ctx sdk.Context, depositorAddress, sportEventUID string, participantID uint64) (val types.Deposit, found bool) {
	sportEventsStore := k.getDepositsStore(ctx)
	depoistKey := types.GetDepositKey(depositorAddress, sportEventUID, participantID)
	b := sportEventsStore.Get(depoistKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetAllDeposits returns all deposits used during genesis dump.
func (k Keeper) GetAllDeposits(ctx sdk.Context) (list []types.Deposit, err error) {
	store := k.getDepositsStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Deposit
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Deposit performs a deposit, set/update everything necessary within the store.
func (k Keeper) Deposit(ctx sdk.Context, creator string, sportEventUID string, amount sdk.Int) (participantID uint64, err error) {
	// Create the deposit object
	deposit := types.NewDeposit(creator, sportEventUID, amount, sdk.ZeroInt(), 0)

	deposit.SetHouseParticipationFee(k.GetHouseParticipationFee(ctx))

	bettorAddress, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		err = sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
		return
	}

	participantID, err = k.orderBookKeeper.AddBookParticipant(
		ctx, bettorAddress, sportEventUID, deposit.Liquidity, deposit.Fee, types.HouseParticipationFeeName,
	)
	if err != nil {
		err = sdkerrors.Wrapf(types.ErrOrderBookDepositProcessing, "%s", err)
		return
	}

	deposit.ParticipantID = participantID

	k.SetDeposit(ctx, deposit)

	return
}
