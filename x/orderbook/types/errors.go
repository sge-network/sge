package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/orderbook module sentinel errors
var (
	ErrOrderBookNotFound                    = sdkerrors.Register(ModuleName, 1001, "order book not found")
	ErrOrderBookNotActive                   = sdkerrors.Register(ModuleName, 1002, "order book not active")
	ErrMaxNumberOfParticipationsReached     = sdkerrors.Register(ModuleName, 1003, "maximum number of participations reached")
	ErrInsufficientUserBalance              = sdkerrors.Register(ModuleName, 1004, "Insufficient user balance.")
	ErrFromBankModule                       = sdkerrors.Register(ModuleName, 1005, "Error returned from Bank Module")
	ErrBookParticipationAlreadyExists       = sdkerrors.Register(ModuleName, 1006, "internal error in setting book participation")
	ErrOrderBookAlreadyPresent              = sdkerrors.Register(ModuleName, 1007, "order book already present")
	ErrDuplicateSenderAndRecipientModule    = sdkerrors.Register(ModuleName, 1008, "sender and receiver module names must not be same")
	ErrInsufficientBalanceInModuleAccount   = sdkerrors.Register(ModuleName, 1009, "Insufficient Balance in Module Account")
	ErrOrderBookParticipationAlreadyPresent = sdkerrors.Register(ModuleName, 1010, "order book participation already present")
	ErrLockAlreadyExists                    = sdkerrors.Register(ModuleName, 1011, "Conflict, lock already exists")
	ErrOrderBookExposureNotFound            = sdkerrors.Register(ModuleName, 1012, "order book exposure not found")
	ErrInsufficientLiquidityInBook          = sdkerrors.Register(ModuleName, 1013, "insufficient liquidity in order book")
	ErrBookParticipationsNotFound           = sdkerrors.Register(ModuleName, 1014, "book participations not found")
	ErrParticipationExposuresNotFound       = sdkerrors.Register(ModuleName, 1015, "participation exposures not found")
	ErrBookParticipationNotFound            = sdkerrors.Register(ModuleName, 1016, "book participation not found")
	ErrParticipationExposureNotFound        = sdkerrors.Register(ModuleName, 1017, "participation exposure not found")
	ErrParticipationExposureAlreadyFilled   = sdkerrors.Register(ModuleName, 1018, "participation exposure already filled")
	ErrInternalProcessingBet                = sdkerrors.Register(ModuleName, 1019, "internal error in processing bet")
	ErrPayoutLockDoesnotExist               = sdkerrors.Register(ModuleName, 1020, "Payout lock for bet uid %s does not exist")
	ErrBookParticipationAlreadySettled      = sdkerrors.Register(ModuleName, 1021, "book participation already settled")
	ErrMismatchInDepositorAddress           = sdkerrors.Register(ModuleName, 1022, "mismatch in depositor address")
	ErrDepositorIsModuleAccount             = sdkerrors.Register(ModuleName, 1023, "depositor is module account")
	ErrWithdrawalAmountIsTooLarge           = sdkerrors.Register(ModuleName, 1024, "withdrawal amount more than available amount for withdrawal")
	ErrMaxWithdrawableAmountIsZero          = sdkerrors.Register(ModuleName, 1025, "maximum withdrawable amount is zero")
)

// x/orderbook module sentinel error text
const (
	ErrTextInvalidDesositor = "invalid depositor address (%s)"
)
