package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrUnlockTokenTimeExpired is the error for unlock time is expired
	ErrUnlockTokenTimeExpired = sdkerrors.Register(ModuleName, 1, "unlock time is expired")

	// ErrSubaccountAlreadyExist is the error for account has already sub account
	ErrSubaccountAlreadyExist = sdkerrors.Register(ModuleName, 2, "account has already sub account")

	// ErrSubaccountDoesNotExist is the error when sub account does not exist
	ErrSubaccountDoesNotExist = sdkerrors.Register(ModuleName, 3, "sub account does not exist")

	// ErrNothingToWithdraw is the error returned when there is nothing to withdraw
	ErrNothingToWithdraw = sdkerrors.Register(ModuleName, 4, "nothing to withdraw")
)
