package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/orderbook/types"
)

type iModuleFunder interface {
	GetModuleAcc() string
}

func (k *Keeper) fund(
	mf iModuleFunder,
	ctx sdk.Context,
	senderAcc sdk.AccAddress,
	amount sdkmath.Int,
) error {
	mAcc := mf.GetModuleAcc()

	// Get the spendable balance of the account holder
	usgeCoins := k.bankKeeper.SpendableCoins(ctx, senderAcc).AmountOf(params.BaseCoinUnit)

	// If account holder has insufficient balance, return error
	if usgeCoins.LT(amount) {
		return sdkerrors.Wrapf(
			types.ErrInsufficientAccountBalance,
			"account Address: %s",
			senderAcc.String(),
		)
	}

	amt := sdk.NewCoins(sdk.NewCoin(params.BaseCoinUnit, amount))

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, mAcc, amt)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrFromBankModule, ": %s", err)
	}

	return nil
}

func (k *Keeper) refund(
	mf iModuleFunder,
	ctx sdk.Context,
	receiverAcc sdk.AccAddress,
	amount sdkmath.Int,
) error {
	mAcc := mf.GetModuleAcc()
	// Get the balance of the sender module account
	balance := k.bankKeeper.GetBalance(
		ctx,
		k.accountKeeper.GetModuleAddress(mAcc),
		params.DefaultBondDenom,
	)

	amt := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount))

	// If module account has insufficient balance, return error
	if balance.Amount.LT(amt.AmountOf(params.DefaultBondDenom)) {
		return sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount,
			mAcc)
	}

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, mAcc,
		receiverAcc, amt)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrFromBankModule, err.Error())
	}

	return nil
}
