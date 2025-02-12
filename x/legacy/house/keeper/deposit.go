package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/legacy/house/types"
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

// GetAllDeposits returns list of deposits.
func (k Keeper) GetAllDeposits(ctx sdk.Context) (list []types.Deposit, err error) {
	store := k.getDepositStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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
	marketUID string, amount sdkmath.Int,
) (participationIndex uint64, feeAmount sdkmath.Int, err error) {
	// Create the deposit object
	deposit := types.NewDeposit(creator, depositor, marketUID, amount, sdkmath.ZeroInt(), 0)

	feeAmount = deposit.CalcHouseParticipationFeeAmount(k.GetHouseParticipationFee(ctx))

	depositorAddr, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		err = sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
		return
	}

	participationIndex, err = k.orderbookKeeper.InitiateOrderBookParticipation(
		ctx, depositorAddr, marketUID, deposit.Amount, feeAmount,
	)
	if err != nil {
		err = sdkerrors.Wrapf(types.ErrOBDepositProcessing, "%s", err)
		return
	}

	deposit.ParticipationIndex = participationIndex

	k.SetDeposit(ctx, deposit)

	return participationIndex, feeAmount, err
}
