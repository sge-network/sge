package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/x/house/types"
)

// SetDeposit sets a deposit.
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	depoistKey := types.GetDepositKey(deposit.Creator, deposit.SportEventUID, deposit.ParticipationIndex)

	store := k.getDepositsStore(ctx)
	b := k.cdc.MustMarshal(&deposit)
	store.Set(depoistKey, b)
}

// GetDeposit returns a specific deposit.
func (k Keeper) GetDeposit(ctx sdk.Context, depositorAddress, sportEventUID string, participationIndex uint64) (val types.Deposit, found bool) {
	sportEventsStore := k.getDepositsStore(ctx)
	depoistKey := types.GetDepositKey(depositorAddress, sportEventUID, participationIndex)
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
func (k Keeper) Deposit(ctx sdk.Context, creator string, sportEventUID string, amount sdk.Int) (participationIndex uint64, err error) {
	// Create the deposit object
	deposit := types.NewDeposit(creator, sportEventUID, amount, sdk.ZeroInt(), 0)

	deposit.SetHouseParticipationFee(k.GetHouseParticipationFee(ctx))

	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		err = sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
		return
	}

	participationIndex, err = k.orderBookKeeper.InitiateBookParticipation(
		ctx, creatorAddr, sportEventUID, deposit.Liquidity, deposit.Fee,
	)
	if err != nil {
		err = sdkerrors.Wrapf(types.ErrOrderBookDepositProcessing, "%s", err)
		return
	}

	deposit.ParticipationIndex = participationIndex

	k.SetDeposit(ctx, deposit)
	emitTransactionEvent(ctx, types.TypeMsgDeposit, cast.ToString(participationIndex), creator)

	return participationIndex, err
}

func emitTransactionEvent(ctx sdk.Context, emitType string, particpationIndex, creator string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			emitType,
			sdk.NewAttribute(types.AttributeKeyParticipationIndex, particpationIndex),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, emitType),
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
		),
	})
}
