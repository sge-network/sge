package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/orderbook module sentinel errors
// nolint
var (
	ErrOrderBookNotFound                  = sdkerrors.Register(ModuleName, 1001, "order book not found")
	ErrOrderBookNotActive                 = sdkerrors.Register(ModuleName, 1002, "order book not active")
	ErrMaxNumberOfParticipantsReached     = sdkerrors.Register(ModuleName, 1003, "maximum number of participants reached")
	ErrInsufficientUserBalance            = sdkerrors.Register(ModuleName, 1004, "Insufficient user balance.")
	ErrFromBankModule                     = sdkerrors.Register(ModuleName, 1005, "Error returned from Bank Module")
	ErrBookParticipantAlreadyExists       = sdkerrors.Register(ModuleName, 1006, "internal error in setting book participant")
	ErrOrderBookAlreadyPresent            = sdkerrors.Register(ModuleName, 1007, "order book already present")
	ErrDuplicateSenderAndRecipientModule  = sdkerrors.Register(ModuleName, 1008, "sender and receiver module names must not be same")
	ErrInsufficientBalanceInModuleAccount = sdkerrors.Register(ModuleName, 1009, "Insufficient Balance in Module Account")
	ErrOrderBookParticipantAlreadyPresent = sdkerrors.Register(ModuleName, 1010, "order book participant already present")
	ErrLockAlreadyExists                  = sdkerrors.Register(ModuleName, 1011, "Conflict, lock already exists")
	ErrOrderBookExposureNotFound          = sdkerrors.Register(ModuleName, 1012, "order book exposure not found")
	ErrInsufficientLiquidityInBook        = sdkerrors.Register(ModuleName, 1013, "insufficient liquidity in order book")
	ErrBookParticipantsNotFound           = sdkerrors.Register(ModuleName, 1014, "book participants not found")
	ErrParticipantExposuresNotFound       = sdkerrors.Register(ModuleName, 1015, "participant exposures not found")
	ErrBookParticipantNotFound            = sdkerrors.Register(ModuleName, 1016, "book participant not found")
	ErrParticipantExposureNotFound        = sdkerrors.Register(ModuleName, 1017, "participant exposure not found")
	ErrParticipantExposureAlreadyFilled   = sdkerrors.Register(ModuleName, 1018, "participant exposure already filled")
	ErrInternalProcessingBet              = sdkerrors.Register(ModuleName, 1019, "internal error in processing bet")
)
