package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/orderbook module sentinel errors
var (
	ErrOrderBookNotFound                    = sdkerrors.Register(ModuleName, 6001, "order book not found")
	ErrOrderBookNotActive                   = sdkerrors.Register(ModuleName, 6002, "order book not active")
	ErrMaxNumberOfParticipationsReached     = sdkerrors.Register(ModuleName, 6003, "maximum number of participations reached")
	ErrInsufficientUserBalance              = sdkerrors.Register(ModuleName, 6004, "Insufficient user balance.")
	ErrFromBankModule                       = sdkerrors.Register(ModuleName, 6005, "Error returned from Bank Module")
	ErrBookParticipationAlreadyExists       = sdkerrors.Register(ModuleName, 6006, "internal error in setting book participation")
	ErrOrderBookAlreadyPresent              = sdkerrors.Register(ModuleName, 6007, "order book already present")
	ErrDuplicateSenderAndRecipientModule    = sdkerrors.Register(ModuleName, 6008, "sender and receiver module names must not be same")
	ErrInsufficientBalanceInModuleAccount   = sdkerrors.Register(ModuleName, 6009, "Insufficient Balance in Module Account")
	ErrOrderBookParticipationAlreadyPresent = sdkerrors.Register(ModuleName, 6010, "order book participation already present")
	ErrLockAlreadyExists                    = sdkerrors.Register(ModuleName, 6011, "Conflict, lock already exists")
	ErrOrderBookExposureNotFound            = sdkerrors.Register(ModuleName, 6012, "order book exposure not found")
	ErrInsufficientLiquidityInBook          = sdkerrors.Register(ModuleName, 6013, "insufficient liquidity in order book")
	ErrBookParticipationsNotFound           = sdkerrors.Register(ModuleName, 6014, "book participations not found")
	ErrParticipationExposuresNotFound       = sdkerrors.Register(ModuleName, 6015, "participation exposures not found")
	ErrBookParticipationNotFound            = sdkerrors.Register(ModuleName, 6016, "book participation not found")
	ErrParticipationExposureNotFound        = sdkerrors.Register(ModuleName, 6017, "participation exposure not found")
	ErrParticipationExposureAlreadyFilled   = sdkerrors.Register(ModuleName, 6018, "participation exposure already filled")
	ErrInternalProcessingBet                = sdkerrors.Register(ModuleName, 6019, "internal error in processing bet")
	ErrPayoutLockDoesnotExist               = sdkerrors.Register(ModuleName, 6020, "Payout lock for bet uid %s does not exist")
	ErrBookParticipationAlreadySettled      = sdkerrors.Register(ModuleName, 6021, "book participation already settled")
	ErrMismatchInDepositorAddress           = sdkerrors.Register(ModuleName, 6022, "mismatch in depositor address")
	ErrDepositorIsModuleAccount             = sdkerrors.Register(ModuleName, 6023, "depositor is module account")
	ErrWithdrawalAmountIsTooLarge           = sdkerrors.Register(ModuleName, 6024, "withdrawal amount more than available amount for withdrawal")
	ErrMaxWithdrawableAmountIsZero          = sdkerrors.Register(ModuleName, 6025, "maximum withdrawable amount is zero")
	ErrParticipationOnInactiveSportEvent    = sdkerrors.Register(ModuleName, 6026, "participation is allowed on an active sport-event only")
	ErrSportEventNotFound                   = sdkerrors.Register(ModuleName, 6027, "sport-event not found to initialize participation")
)

// x/orderbook module sentinel error text
const (
	ErrTextInvalidDesositor = "invalid depositor address (%s)"
)
