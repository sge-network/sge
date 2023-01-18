package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"
)

// x/bet module sentinel errors
// nolint
var (
	ErrInvalidBetUID                        = sdkerrors.Register(ModuleName, 2001, "bet UID is invalid")
	ErrInVerification                       = sdkerrors.Register(ModuleName, 2002, "error in verification process")
	ErrBetIsCanceled                        = sdkerrors.Register(ModuleName, 2003, "bet is already canceled")
	ErrBetIsSettled                         = sdkerrors.Register(ModuleName, 2004, "bet is already settled")
	ErrInSRBettorLoses                      = sdkerrors.Register(ModuleName, 2005, "internal error in processing loss in SR")
	ErrInSRBettorWins                       = sdkerrors.Register(ModuleName, 2006, "internal error in processing win in SR")
	ErrNoMatchingSportEvent                 = sdkerrors.Register(ModuleName, 2007, "sport event is not found")
	ErrInvalidOddsUID                       = sdkerrors.Register(ModuleName, 2008, "valid odds UID should be provided")
	ErrEmptyOddsValue                       = sdkerrors.Register(ModuleName, 2009, "odds value should not be empty")
	ErrInvalidSportEventUID                 = sdkerrors.Register(ModuleName, 2010, "valid sport event UID should be provided")
	ErrInvalidTicket                        = sdkerrors.Register(ModuleName, 2011, "valid ticket should be provided")
	ErrInvalidAmount                        = sdkerrors.Register(ModuleName, 2012, "valid amount should be provided")
	ErrNoMatchingBet                        = sdkerrors.Register(ModuleName, 2013, "no matching bet")
	ErrResultNotDeclared                    = sdkerrors.Register(ModuleName, 2014, "sport event result is not declared")
	ErrDuplicateUID                         = sdkerrors.Register(ModuleName, 2015, "UID is already set")
	ErrInSRPlacementProcessing              = sdkerrors.Register(ModuleName, 2016, "internal error in processing bet placement in SR")
	ErrSportEventStatusNotPending           = sdkerrors.Register(ModuleName, 2017, "can not place bet on this sport event any more")
	ErrEndTSIsPassed                        = sdkerrors.Register(ModuleName, 2018, "the sport event is expired")
	ErrOddsUIDNotExist                      = sdkerrors.Register(ModuleName, 2019, "sport event does not have this odds UID")
	ErrInSRRefund                           = sdkerrors.Register(ModuleName, 2020, "internal error in refunding user in SR")
	ErrInactiveSportEvent                   = sdkerrors.Register(ModuleName, 2021, "the sport event is not active")
	ErrBetAmountIsLow                       = sdkerrors.Register(ModuleName, 2022, "bet amount is lower than the minimum allowed")
	ErrInConvertingOddsToDec                = sdkerrors.Register(ModuleName, 2023, "internal error in converting odds value from string to sdk.Dec")
	ErrInConvertingOddsToInt                = sdkerrors.Register(ModuleName, 2024, "internal error in converting odds value from string to sdk.Int")
	ErrOddsDataNotFound                     = sdkerrors.Register(ModuleName, 2025, "odds does not exist in ticket payload")
	ErrInvalidOddsType                      = sdkerrors.Register(ModuleName, 2026, "valid odds type should be provided, 1: decimal, 2: fractional, 3: monyline")
	ErrUserKycFailed                        = sdkerrors.Register(ModuleName, 2027, "the bettor failed the KYC Validation")
	ErrNoKycField                           = sdkerrors.Register(ModuleName, 2028, "KYC field does not exist in ticket payload")
	ErrNoKycIdField                         = sdkerrors.Register(ModuleName, 2029, "KYC ID does not exist in KYC part of ticket payload")
	ErrCanNotQueryLargeNumberOfBets         = sdkerrors.Register(ModuleName, 2030, "can not query more than "+cast.ToString(MaxAllowedQueryBetsCount))
	ErrDecimalOddsCanNotBeNegative          = sdkerrors.Register(ModuleName, 2031, "decimal odds value can not bet negative")
	ErrDecimalOddsCanNotBeLessThanOne       = sdkerrors.Register(ModuleName, 2032, "decimal odds value can not less then 1")
	ErrFractionalOddsCanNotBeNegativeOrZero = sdkerrors.Register(ModuleName, 2033, "fractional odds numbers can not be negative")
	ErrMoneylineOddsCanNotBeZero            = sdkerrors.Register(ModuleName, 2034, "moneyline odds can not be zero")
	ErrFractionalOddsIncorrectFormat        = sdkerrors.Register(ModuleName, 2035, "incorrect format of fractional odds")
	ErrBettorAddressNotEqualToCreator       = sdkerrors.Register(ModuleName, 2036, "provided bettor address is not equal to bet owner")
)

// x/bet module sentinel error text
// nolint
const (
	ErrTextInvalidCreator                         = "invalid creator address (%s)"
	ErrTextInitGenesisFailedBecauseOfMissingBetID = "no bet id found for the bet with uuid (%s)"
)
