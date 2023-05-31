package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/strategicreserve module sentinel errors
var (
	ErrInTicketVerification                     = sdkerrors.Register(ModuleName, 6001, "ticket verification failed")
	ErrInTicketPayloadValidation                = sdkerrors.Register(ModuleName, 6002, "ticket validation failed")
	ErrOrderBookNotFound                        = sdkerrors.Register(ModuleName, 6003, "order book not found")
	ErrOrderBookNotActive                       = sdkerrors.Register(ModuleName, 6004, "order book not active")
	ErrMaxNumberOfParticipationsReached         = sdkerrors.Register(ModuleName, 6005, "maximum number of participations reached")
	ErrInsufficientAccountBalance               = sdkerrors.Register(ModuleName, 6006, "insufficient account balance.")
	ErrFromBankModule                           = sdkerrors.Register(ModuleName, 6007, "error returned from Bank Module")
	ErrBookParticipationAlreadyExists           = sdkerrors.Register(ModuleName, 6008, "internal error in setting book participation")
	ErrOrderBookAlreadyPresent                  = sdkerrors.Register(ModuleName, 6009, "order book already present")
	ErrDuplicateSenderAndRecipientModule        = sdkerrors.Register(ModuleName, 6010, "sender and receiver module names must not be same")
	ErrInsufficientBalanceInModuleAccount       = sdkerrors.Register(ModuleName, 6011, "Insufficient Balance in Module Account")
	ErrOrderBookParticipationAlreadyPresent     = sdkerrors.Register(ModuleName, 6012, "order book participation already present")
	ErrLockAlreadyExists                        = sdkerrors.Register(ModuleName, 6013, "Conflict, lock already exists")
	ErrOrderBookExposureNotFound                = sdkerrors.Register(ModuleName, 6014, "order book exposure not found")
	ErrInsufficientLiquidityInOrderBook         = sdkerrors.Register(ModuleName, 6015, "insufficient liquidity in order book")
	ErrBookParticipationsNotFound               = sdkerrors.Register(ModuleName, 6016, "book participations not found")
	ErrParticipationExposuresNotFound           = sdkerrors.Register(ModuleName, 6017, "participation exposures not found")
	ErrOrderBookParticipationNotFound           = sdkerrors.Register(ModuleName, 6018, "book participation not found")
	ErrParticipationExposureNotFound            = sdkerrors.Register(ModuleName, 6019, "participation exposure not found")
	ErrParticipationExposureAlreadyFilled       = sdkerrors.Register(ModuleName, 6020, "participation exposure already filled")
	ErrInternalProcessingBet                    = sdkerrors.Register(ModuleName, 6021, "internal error in processing bet")
	ErrPayoutLockDoesnotExist                   = sdkerrors.Register(ModuleName, 6022, "Payout lock for bet uid %s does not exist")
	ErrBookParticipationAlreadySettled          = sdkerrors.Register(ModuleName, 6023, "book participation already settled")
	ErrMismatchInDepositorAddress               = sdkerrors.Register(ModuleName, 6024, "mismatch in depositor address")
	ErrWithdrawalAmountIsTooLarge               = sdkerrors.Register(ModuleName, 6025, "withdrawal amount more than available amount for withdrawal")
	ErrMaxWithdrawableAmountIsZero              = sdkerrors.Register(ModuleName, 6026, "maximum withdrawable amount is zero")
	ErrParticipationOnInactiveMarket            = sdkerrors.Register(ModuleName, 6027, "participation is allowed on an active market only")
	ErrMarketNotFound                           = sdkerrors.Register(ModuleName, 6028, "market not found to initialize participation")
	ErrTranferringDepositorProfit               = sdkerrors.Register(ModuleName, 6029, "error while transferring the loser bat amount and profit to depositor")
	ErrDepositNotFoundForParticipation          = sdkerrors.Register(ModuleName, 6030, "corresponding deposit not found for the participation")
	ErrInvalidDataFeeCollectorProposalAmount    = sdkerrors.Register(ModuleName, 6031, "invalid data fee collector feed proposal amount")
	ErrFeeGrantExists                           = sdkerrors.Register(ModuleName, 6032, "fee grant already exists for the grantee")
	ErrSDKFeeGrantExists                        = sdkerrors.Register(ModuleName, 6033, "sdk fee grant already exists for the grantee")
	ErrInFeeGrantAllowance                      = sdkerrors.Register(ModuleName, 6034, "fee grant allowance grant failed")
	ErrInSrPoolFeeGrant                         = sdkerrors.Register(ModuleName, 6035, "sr pool sdk fee grant failed")
	ErrParticipationsCanNotCoverthePayoutProfit = sdkerrors.Register(ModuleName, 6036, "not enought fund in the participations to cover the payout")
)

// x/strategicreserve module sentinel error text
const (
	ErrTextInvalidDesositor = "invalid depositor address (%s)"
)
