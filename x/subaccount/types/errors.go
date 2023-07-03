package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrUnlockTokenTimeExpired is the error for unlock time is expired
	ErrUnlockTokenTimeExpired = sdkerrors.Register(ModuleName, 1, "token ")
)
