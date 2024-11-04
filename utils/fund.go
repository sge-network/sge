package utils

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
)

type IModuleFunder interface {
	GetModuleAcc() string
}

type ModuleAccFunder struct {
	bk BankKeeper
	ak AccountKeeper

	bankError error
}

func NewModuleAccFunder(bk BankKeeper, ak AccountKeeper, bankError error) *ModuleAccFunder {
	return &ModuleAccFunder{bk, ak, bankError}
}

// AccountKeeper defines the expected account keeper methods.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper methods.
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, ecipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// Fund transfers the input amount from sender to the module account.
func (f *ModuleAccFunder) Fund(
	mf IModuleFunder,
	ctx sdk.Context,
	senderAcc sdk.AccAddress,
	amount sdkmath.Int,
) error {
	amt := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount))

	// Transfer funds
	err := f.bk.SendCoinsFromAccountToModule(ctx, senderAcc, mf.GetModuleAcc(), amt)
	if err != nil {
		return sdkerrors.Wrapf(f.bankError, ": %s", err)
	}

	return nil
}

// Refund transfers the input amount from the module account to the receiver account.
func (f *ModuleAccFunder) Refund(
	mf IModuleFunder,
	ctx sdk.Context,
	receiverAcc sdk.AccAddress,
	amount sdkmath.Int,
) error {
	mAcc := mf.GetModuleAcc()

	amt := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount))

	// Transfer funds
	err := f.bk.SendCoinsFromModuleToAccount(ctx, mAcc, receiverAcc, amt)
	if err != nil {
		return sdkerrors.Wrap(f.bankError, err.Error())
	}

	return nil
}
