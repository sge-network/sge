package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/strategicreserve module sentinel errors
// nolint
var (
	ErrInsufficientUserBalance            = sdkerrors.Register(ModuleName, 1100, "Insufficient user balance. User Address: %s")
	ErrInsufficientUnlockedAmountInSrPool = sdkerrors.Register(ModuleName, 1101, "Unlocked amount in the SR_Pool is insufficient")
	ErrInsufficientLockedAmountInSrPool   = sdkerrors.Register(ModuleName, 1102, "Locked amount in the SR_Pool is insufficient")
	ErrInsufficientBalanceInModuleAccount = sdkerrors.Register(ModuleName, 1103, "Balance in the %s Module Account is insufficient")
	ErrFromBankModule                     = sdkerrors.Register(ModuleName, 1104, "Error returned from Bank Module: %s")
	ErrPayoutLockDoesnotExist             = sdkerrors.Register(ModuleName, 1105, "Payout lock for bet uid %s does not exist")
	ErrLockAlreadyExists                  = sdkerrors.Register(ModuleName, 1106, "Conflict, lock already exists")
	ErrDuplicateSenderAndRecipientModule  = sdkerrors.Register(ModuleName, 1114, "sender and receiver module names cannot be same")
)
