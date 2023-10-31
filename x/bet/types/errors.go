package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bet module sentinel errors
var (
	ErrInTicketVerification                 = sdkerrors.Register(ModuleName, 2001, "ticket verification failed")
	ErrInTicketValidation                   = sdkerrors.Register(ModuleName, 2002, "ticket validation failed")
	ErrInWager                              = sdkerrors.Register(ModuleName, 2003, "bet placement failed")
	ErrInvalidBetUID                        = sdkerrors.Register(ModuleName, 2004, "invalid Bet Uid")
	ErrBetIsCanceled                        = sdkerrors.Register(ModuleName, 2005, "bet is already canceled")
	ErrBetIsSettled                         = sdkerrors.Register(ModuleName, 2006, "bet is already settled")
	ErrInOBBettorLoses                      = sdkerrors.Register(ModuleName, 2007, "internal error in processing better loss in orderbook")
	ErrInOBBettorWins                       = sdkerrors.Register(ModuleName, 2008, "internal error in processing better win in orderbook")
	ErrNoMatchingMarket                     = sdkerrors.Register(ModuleName, 2009, "market not found")
	ErrInvalidOddsUID                       = sdkerrors.Register(ModuleName, 2010, "invalid odds UID")
	ErrEmptyOddsValue                       = sdkerrors.Register(ModuleName, 2011, "odds value should not be empty")
	ErrInvalidMarketUID                     = sdkerrors.Register(ModuleName, 2012, "invalid market UID")
	ErrInvalidTicket                        = sdkerrors.Register(ModuleName, 2013, "invalid ticket")
	ErrInvalidAmount                        = sdkerrors.Register(ModuleName, 2014, "invalid amount")
	ErrNoMatchingBet                        = sdkerrors.Register(ModuleName, 2015, "bet not found")
	ErrResultNotDeclared                    = sdkerrors.Register(ModuleName, 2016, "market result is not declared")
	ErrDuplicateUID                         = sdkerrors.Register(ModuleName, 2017, "UID is already set")
	ErrInOBWagerProcessing                  = sdkerrors.Register(ModuleName, 2018, "internal error in processing bet placement in orderbook")
	ErrEndTSIsPassed                        = sdkerrors.Register(ModuleName, 2019, "market is expired")
	ErrOddsUIDNotExist                      = sdkerrors.Register(ModuleName, 2020, "market does not have this odds UID")
	ErrInOBRefund                           = sdkerrors.Register(ModuleName, 2021, "internal error in refunding user in Orderbook")
	ErrInactiveMarket                       = sdkerrors.Register(ModuleName, 2022, "market is not active")
	ErrBetAmountIsLow                       = sdkerrors.Register(ModuleName, 2023, "bet amount is lower than the minimum allowed")
	ErrDecimalOddsIncorrectFormat           = sdkerrors.Register(ModuleName, 2024, "incorrect format of decimal odds value")
	ErrFractionalOddsIncorrectFormat        = sdkerrors.Register(ModuleName, 2025, "incorrect format of fractional odds value")
	ErrMoneylineOddsIncorrectFormat         = sdkerrors.Register(ModuleName, 2026, "incorrect format of moneyline odds value")
	ErrOddsDataNotFound                     = sdkerrors.Register(ModuleName, 2027, "odds does not exist in ticket payload")
	ErrUserKycFailed                        = sdkerrors.Register(ModuleName, 2028, "the bettor failed the KYC Validation")
	ErrCanNotQueryLargeNumberOfBets         = sdkerrors.Register(ModuleName, 2029, "large amount of bets requested")
	ErrDecimalOddsShouldBePositive          = sdkerrors.Register(ModuleName, 2030, "decimal odds value should be positive")
	ErrDecimalOddsCanNotBeLessThanOne       = sdkerrors.Register(ModuleName, 2032, "decimal odds value can not less than or equal to 1")
	ErrBettorAddressNotEqualToCreator       = sdkerrors.Register(ModuleName, 2033, "provided bettor address is not equal to bet owner")
	ErrMaxLossMultiplierCanNotBeZero        = sdkerrors.Register(ModuleName, 2034, "max loss multiplier cannot be nil or zero")
	ErrMaxLossMultiplierCanNotBeMoreThanOne = sdkerrors.Register(ModuleName, 2035, "max loss multiplier cannot be more than one")
	ErrInsufficientOdds                     = sdkerrors.Register(ModuleName, 2036, "market odds length not same as odds sent in wager")
)

// x/bet module sentinel error text
const (
	ErrTextInvalidParamType                                  = "invalid parameter type"
	ErrTextBatchSettlementCountMustBePositive                = "batch settlement count should be a positive number"
	ErrTextMaxBetUIDQueryCountMustBePositive                 = "max bet by uid query count should be a positive number"
	ErrTextInitGenesisFailedBecauseOfMissingBetID            = "no bet id found for the bet with uuid"
	ErrTextInitGenesisFailedBecauseOfNotEqualStats           = "bet list items count is not equal to stats count"
	ErrTextInitGenesisFailedBetCountNotEqualActiveAndSettled = "sum of active and settled list items count is not equal to bet list items count"
	ErrTextInitGenesisFailedNotActiveOrSettled               = "bet is not active nor settled with uuid"
	ErrTextInitGenesisFailedSettlementHeightIsZero           = "settlement height can not be zero for the settled bet with uuid"
	ErrTextInitGenesisFailedSettlementHeightIsZeroForList    = "settlement height can not be zero for a bet in the settled bet list with uuid"
	ErrTextInitGenesisFailedSettlementHeightIsNotZero        = "settlement height should be zero for the active bet with uuid"
)
