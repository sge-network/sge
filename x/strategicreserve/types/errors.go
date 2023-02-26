package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/strategicreserve module sentinel errors
var (
	ErrInsufficientUserBalance            = sdkerrors.Register(ModuleName, 7001, "Insufficient user balance. User Address: %s")
	ErrInsufficientUnlockedAmountInSrPool = sdkerrors.Register(ModuleName, 7002, "Insufficient funds in SR")
	ErrInsufficientLockedAmountInSrPool   = sdkerrors.Register(ModuleName, 7003, "Insufficient funds locked in SR Pool")
	ErrInsufficientBalanceInModuleAccount = sdkerrors.Register(ModuleName, 7004, "Insufficient Balance in the %s Module Account")
	ErrFromBankModule                     = sdkerrors.Register(ModuleName, 7005, "Error returned from Bank Module: %s")
	ErrPayoutLockDoesnotExist             = sdkerrors.Register(ModuleName, 7006, "Payout lock for bet uid %s does not exist")
	ErrLockAlreadyExists                  = sdkerrors.Register(ModuleName, 7007, "Conflict, lock already exists")
	ErrDuplicateSenderAndRecipientModule  = sdkerrors.Register(ModuleName, 7008, "sender and receiver module names must not be same")
	ErrTextNilReserver                    = sdkerrors.Register(ModuleName, 7009, "Reserver must not be nil")
)
