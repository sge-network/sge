package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// transferFundsFromAccountToModule transfers the given amount from
// the given account address to the module account passed.
// Returns an error if the account holder has insufficient balance.
func (k Keeper) transferFundsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, moduleAccName string, amount sdk.Int) error {
	// Get the spendable balance of the account holder
	usgeCoins := k.bankKeeper.SpendableCoins(ctx, address).AmountOf(params.BaseCoinUnit)

	// If account holder has insufficient balance, return error
	if usgeCoins.LT(amount) {
		return sdkerrors.Wrapf(types.ErrInsufficientAccountBalance, "account Address: %s", address.String())
	}

	amt := sdk.NewCoins(sdk.NewCoin(params.BaseCoinUnit, amount))

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, moduleAccName, amt)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrFromBankModule, ": %s", err)
	}

	return nil
}

// transferFundsFromModuleToModule transfers the given amount from a module
// account to another module account.
// Returns an error if the sender module has insufficient balance.
func (k Keeper) transferFundsFromModuleToModule(
	ctx sdk.Context,
	senderModule string,
	recipientModule string,
	amount sdk.Int,
) error {
	if senderModule == recipientModule {
		return types.ErrDuplicateSenderAndRecipientModule
	}

	amt := sdk.NewCoins(sdk.NewCoin(params.BaseCoinUnit, amount))

	// Get the balance of the sender module account
	balance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(senderModule), params.BaseCoinUnit)

	// If sender module account has insufficient balance, return error
	if balance.Amount.LT(amt.AmountOf(params.BaseCoinUnit)) {
		return sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount, "%s", senderModule)
	}

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrFromBankModule, err.Error())
	}

	return nil
}

// transferFundsFromModuleToAccount transfers the given amount from a module
// account to the given account address.
// Returns an error if the account holder has insufficient balance.
func (k Keeper) transferFundsFromModuleToAccount(ctx sdk.Context,
	moduleAccName string, address sdk.AccAddress, amount sdk.Int,
) error {
	// Get the balance of the sender module account
	balance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(
		moduleAccName), params.DefaultBondDenom)

	amt := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount))

	// If module account has insufficient balance, return error
	if balance.Amount.LT(amt.AmountOf(params.DefaultBondDenom)) {
		return sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount,
			moduleAccName)
	}

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccName,
		address, amt)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrFromBankModule, err.Error())
	}

	return nil
}
