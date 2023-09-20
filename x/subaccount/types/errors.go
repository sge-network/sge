package types

// DONTCOVER

import (
	cosmerrors "cosmossdk.io/errors"
)

var (
	ErrUnlockTokenTimeExpired = cosmerrors.Register(ModuleName, 1, "unlock time is expired")
	ErrSubaccountAlreadyExist = cosmerrors.Register(ModuleName, 2, "account has already sub account")
	ErrSubaccountDoesNotExist = cosmerrors.Register(ModuleName, 3, "sub account does not exist")
	ErrNothingToWithdraw      = cosmerrors.Register(ModuleName, 4, "nothing to withdraw")
	ErrInvalidLockedBalance   = cosmerrors.Register(ModuleName, 5, "invalid locked balance")
)
