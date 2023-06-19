package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/orderbook module sentinel errors
var (
	ErrInTicketVerification               = sdkerrors.Register(ModuleName, 6001, "ticket verification failed")
	ErrInTicketPayloadValidation          = sdkerrors.Register(ModuleName, 6002, "ticket validation failed")
	ErrOrderBookNotFound                  = sdkerrors.Register(ModuleName, 6003, "order book not found")
	ErrOrderBookNotActive                 = sdkerrors.Register(ModuleName, 6004, "order book not active")
	ErrMaxNumberOfParticipationsReached   = sdkerrors.Register(ModuleName, 6005, "maximum number of participations reached")
	ErrInsufficientAccountBalance         = sdkerrors.Register(ModuleName, 6006, "insufficient account balance.")
	ErrFromBankModule                     = sdkerrors.Register(ModuleName, 6007, "error returned from Bank Module")
	ErrBookParticipationAlreadyExists     = sdkerrors.Register(ModuleName, 6008, "internal error in setting book participation")
	ErrOrderBookAlreadyPresent            = sdkerrors.Register(ModuleName, 6009, "order book already present")
	ErrInsufficientBalanceInModuleAccount = sdkerrors.Register(ModuleName, 6010, "Insufficient Balance in Module Account")
	ErrOrderBookExposureNotFound          = sdkerrors.Register(ModuleName, 6011, "order book exposure not found")
	ErrBookParticipationsNotFound         = sdkerrors.Register(ModuleName, 6012, "book participations not found")
	ErrParticipationExposuresNotFound     = sdkerrors.Register(ModuleName, 6013, "participation exposures not found")
	ErrOrderBookParticipationNotFound     = sdkerrors.Register(ModuleName, 6014, "book participation not found")
	ErrParticipationExposureNotFound      = sdkerrors.Register(ModuleName, 6015, "participation exposure not found")
	ErrParticipationExposureAlreadyFilled = sdkerrors.Register(ModuleName, 6016, "participation exposure already filled")
	ErrInternalProcessingBet              = sdkerrors.Register(ModuleName, 6017, "internal error in processing bet")
	ErrBookParticipationAlreadySettled    = sdkerrors.Register(ModuleName, 6018, "book participation already settled")
	ErrMismatchInDepositorAddress         = sdkerrors.Register(ModuleName, 6019, "mismatch in depositor address")
	ErrWithdrawalAmountIsTooLarge         = sdkerrors.Register(ModuleName, 6020, "withdrawal amount more than available amount for withdrawal")
	ErrMaxWithdrawableAmountIsZero        = sdkerrors.Register(ModuleName, 6021, "maximum withdrawal amount is zero")
	ErrParticipationOnInactiveMarket      = sdkerrors.Register(ModuleName, 6022, "participation is allowed on an active market only")
	ErrMarketNotFound                     = sdkerrors.Register(ModuleName, 6023, "market not found to initialize participation")
	ErrInsufficientFundToCoverPayout      = sdkerrors.Register(ModuleName, 6024, "insufficient fund in the participations to cover the payout")
	ErrUnknownMarketStatus                = sdkerrors.Register(ModuleName, 6025, "unknown market status of orderbook settlement")
	ErrWithdrawalTooLarge                 = sdkerrors.Register(ModuleName, 6026, "withdrawal is more than unused amount")
)

// ErrTextInvalidDepositor x/orderbook module sentinel error text
const (
	ErrTextInvalidDepositor = "invalid depositor address (%s)"
)
