package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnlockTokenTimeExpired = sdkerrors.Register(ModuleName, 1, "unlock time is expired")
	ErrSubaccountAlreadyExist = sdkerrors.Register(ModuleName, 2, "account has already sub account")
	ErrSubaccountDoesNotExist = sdkerrors.Register(ModuleName, 3, "sub account does not exist")
	ErrNothingToWithdraw      = sdkerrors.Register(ModuleName, 4, "nothing to withdraw")
)
