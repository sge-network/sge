package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"
)

// x/bet module sentinel errors
var (
	ErrInvalidBetUID                        = sdkerrors.Register(ModuleName, 2001, "bet UID is invalid")
	ErrInVerification                       = sdkerrors.Register(ModuleName, 2002, "error in verification process")
	ErrBetIsCanceled                        = sdkerrors.Register(ModuleName, 2003, "bet is already canceled")
	ErrBetIsSettled                         = sdkerrors.Register(ModuleName, 2004, "bet is already settled")
	ErrInSRBettorLoses                      = sdkerrors.Register(ModuleName, 2005, "internal error in processing loss in SR")
	ErrInSRBettorWins                       = sdkerrors.Register(ModuleName, 2006, "internal error in processing win in SR")
	ErrNoMatchingSportEvent                 = sdkerrors.Register(ModuleName, 2007, "sport-event not found")
	ErrInvalidOddsUID                       = sdkerrors.Register(ModuleName, 2008, "valid odds UID should be provided")
	ErrEmptyOddsValue                       = sdkerrors.Register(ModuleName, 2009, "odds value should not be empty")
	ErrInvalidSportEventUID                 = sdkerrors.Register(ModuleName, 2010, "valid sport-event UID should be provided")
	ErrInvalidTicket                        = sdkerrors.Register(ModuleName, 2011, "valid ticket should be provided")
	ErrInvalidAmount                        = sdkerrors.Register(ModuleName, 2012, "valid amount should be provided")
	ErrNoMatchingBet                        = sdkerrors.Register(ModuleName, 2013, "no matching bet")
	ErrResultNotDeclared                    = sdkerrors.Register(ModuleName, 2014, "sport-event result is not declared")
	ErrDuplicateUID                         = sdkerrors.Register(ModuleName, 2015, "UID is already set")
	ErrInSRPlacementProcessing              = sdkerrors.Register(ModuleName, 2016, "internal error in processing bet placement in SR")
	ErrEndTSIsPassed                        = sdkerrors.Register(ModuleName, 2017, "the sport-event is expired")
	ErrOddsUIDNotExist                      = sdkerrors.Register(ModuleName, 2018, "sport-event does not have this odds UID")
	ErrInSRRefund                           = sdkerrors.Register(ModuleName, 2019, "internal error in refunding user in SR")
	ErrInactiveSportEvent                   = sdkerrors.Register(ModuleName, 2020, "the sport-event is not active")
	ErrBetAmountIsLow                       = sdkerrors.Register(ModuleName, 2021, "bet amount is lower than the minimum allowed")
	ErrInConvertingOddsToDec                = sdkerrors.Register(ModuleName, 2022, "internal error in converting odds value from string to sdk.Dec")
	ErrInConvertingOddsToInt                = sdkerrors.Register(ModuleName, 2023, "internal error in converting odds value from string to sdk.Int")
	ErrOddsDataNotFound                     = sdkerrors.Register(ModuleName, 2024, "odds does not exist in ticket payload")
	ErrInvalidOddsType                      = sdkerrors.Register(ModuleName, 2025, "valid odds type should be provided, 1: decimal, 2: fractional, 3: moneyline")
	ErrUserKycFailed                        = sdkerrors.Register(ModuleName, 2026, "the bettor failed the KYC Validation")
	ErrCanNotQueryLargeNumberOfBets         = sdkerrors.Register(ModuleName, 2027, "can not query more than "+cast.ToString(MaxAllowedQueryBetsCount))
	ErrDecimalOddsShouldBePositive          = sdkerrors.Register(ModuleName, 2028, "decimal odds value should be positive")
	ErrDecimalOddsCanNotBeLessThanOne       = sdkerrors.Register(ModuleName, 2029, "decimal odds value can not less than or equal to 1")
	ErrFractionalOddsCanNotBeNegativeOrZero = sdkerrors.Register(ModuleName, 2030, "fractional odds numbers can not be negative")
	ErrMoneylineOddsCanNotBeZero            = sdkerrors.Register(ModuleName, 2031, "moneyline odds can not be zero")
	ErrFractionalOddsIncorrectFormat        = sdkerrors.Register(ModuleName, 2032, "incorrect format of fractional odds")
	ErrBettorAddressNotEqualToCreator       = sdkerrors.Register(ModuleName, 2033, "provided bettor address is not equal to bet owner")
	ErrMaxLossMultiplierCanNotBeZero        = sdkerrors.Register(ModuleName, 2034, "max loss multiplier cannot be nil or zero")
	ErrMaxLossMultiplierCanNotBeMoreThanOne = sdkerrors.Register(ModuleName, 2035, "max loss multiplier cannot be more than one")
)

// x/bet module sentinel error text
const (
	ErrTextInvalidParamType                                  = "invalid parameter type: %T"
	ErrTextBatchSettlementCountMustBePositive                = "batch settlement count should be a positive number %d"
	ErrTextInvalidCreator                                    = "invalid creator address (%s)"
	ErrTextInitGenesisFailedBecauseOfMissingBetID            = "no bet id found for the bet with uuid (%s)"
	ErrTextInitGenesisFailedBecauseOfNotEqualStats           = "bet list items count (%d) is not equal to stats count (%d)"
	ErrTextInitGenesisFailedBetCountNotEqualActiveAndSettled = "sum of active and settled list items count (%d) is not equal to bet list items count (%d)"
	ErrTextInitGenesisFailedNotActiveOrSettled               = "bet is not active nor settled with uuid (%s)"
	ErrTextInitGenesisFailedSettlementHeightIsZero           = "settlement height can not be zero for the settled bet with uuid (%s)"
	ErrTextInitGenesisFailedSettlementHeightIsZeroForList    = "settlement height can not be zero for a bet in the settled bet list with uuid (%s)"
	ErrTextInitGenesisFailedSettlementHeightIsNotZero        = "settlement height should be zero for the active bet with uuid (%s)"
)
