package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrUnlockTokenTimeExpired is the error for unlock time is expired
	ErrUnlockTokenTimeExpired = sdkerrors.Register(ModuleName, 1, "unlock time is expired")

	// ErrSubaccountAlreadyExist is the error for account has already sub account
	ErrSubaccountAlreadyExist = sdkerrors.Register(ModuleName, 2, "account has already sub account")
)
