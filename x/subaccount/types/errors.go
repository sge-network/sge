package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrInTicketVerification      = sdkerrors.Register(ModuleName, 8001, "error in ticket verification process")
	ErrInTicketPayloadValidation = sdkerrors.Register(ModuleName, 8002, "error in ticket payload validation")
	ErrUnlockTokenTimeExpired    = sdkerrors.Register(ModuleName, 8003, "unlock time is expired")
	ErrSubaccountAlreadyExist    = sdkerrors.Register(ModuleName, 8004, "account has already sub account")
	ErrSubaccountDoesNotExist    = sdkerrors.Register(ModuleName, 8005, "sub account does not exist")
	ErrNothingToWithdraw         = sdkerrors.Register(ModuleName, 8006, "nothing to withdraw")
	ErrInvalidLockedBalance      = sdkerrors.Register(ModuleName, 8007, "invalid locked balance")
	ErrSendCoinError             = sdkerrors.Register(ModuleName, 8008, "send coin error")
	ErrWithdrawLocked            = sdkerrors.Register(ModuleName, 8009, "withdrawal of locked coins failed")
)
