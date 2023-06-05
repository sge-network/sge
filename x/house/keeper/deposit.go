package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/x/house/types"
)

// SetDeposit sets a deposit in the store
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	depositKey := types.GetDepositKey(deposit.DepositorAddress, deposit.MarketUID,
		deposit.ParticipationIndex)

	store := k.getDepositStore(ctx)
	b := k.cdc.MustMarshal(&deposit)
	store.Set(depositKey, b)
}

// GetDeposit returns a specific deposit from the store.
func (k Keeper) GetDeposit(ctx sdk.Context, depositorAddress,
	marketUID string, participationIndex uint64,
) (val types.Deposit, found bool) {
	marketsStore := k.getDepositStore(ctx)
	depositKey := types.GetDepositKey(depositorAddress, marketUID, participationIndex)
	b := marketsStore.Get(depositKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetDeposits returns list of deposits.
func (k Keeper) GetDeposits(ctx sdk.Context) (list []types.Deposit, err error) {
	store := k.getDepositStore(ctx)
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

// Deposit performs a deposit transaction and stores a new deposit in store.
func (k Keeper) Deposit(ctx sdk.Context, creator, depositor string,
	marketUID string, amount sdk.Int,
) (participationIndex uint64, err error) {
	// Create the deposit object
	deposit := types.NewDeposit(creator, depositor, marketUID, amount, sdk.ZeroInt(), 0)

	deposit.SetHouseParticipationFee(k.GetHouseParticipationFee(ctx))

	depositorAddr, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		err = sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
		return
	}

	participationIndex, err = k.srKeeper.InitiateOrderBookParticipation(
		ctx, depositorAddr, marketUID, deposit.Liquidity, deposit.Fee,
	)
	if err != nil {
		err = sdkerrors.Wrapf(types.ErrSRDepositProcessing, "%s", err)
		return
	}

	deposit.ParticipationIndex = participationIndex

	k.SetDeposit(ctx, deposit)
	emitTransactionEvent(
		ctx,
		types.TypeMsgDeposit,
		cast.ToString(participationIndex),
		creator,
		depositor,
	)

	return participationIndex, err
}

func emitTransactionEvent(
	ctx sdk.Context,
	emitType string,
	participationIndex, creator, depositor string,
) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			emitType,
			sdk.NewAttribute(types.AttributeKeyParticipationIndex, participationIndex),
			sdk.NewAttribute(types.AttributeKeyParticipationIndex, depositor),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, emitType),
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
		),
	})
}
