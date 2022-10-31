package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bet module sentinel errors
// nolint
var (
	ErrInvalidBetUID              = sdkerrors.Register(ModuleName, 2001, "bet UID is invalid")
	ErrInVerification             = sdkerrors.Register(ModuleName, 2002, "error in verification process")
	ErrBetIsInvalid               = sdkerrors.Register(ModuleName, 2003, "bet is invalid")
	ErrBetIsAborted               = sdkerrors.Register(ModuleName, 2004, "bet is already aborted")
	ErrBetIsCanceled              = sdkerrors.Register(ModuleName, 2005, "bet is already canceled")
	ErrBetIsSettled               = sdkerrors.Register(ModuleName, 2006, "bet is already settled")
	ErrInUnmarshaling             = sdkerrors.Register(ModuleName, 2007, "internal error in unmarshaling")
	ErrInSRBettorLoses            = sdkerrors.Register(ModuleName, 2008, "internal error in processing loss in SR")
	ErrInSRBettorWins             = sdkerrors.Register(ModuleName, 2009, "internal error in processing win in SR")
	ErrNoMatchingSportEvent       = sdkerrors.Register(ModuleName, 2010, "sport event is not found")
	ErrEmptyBetUID                = sdkerrors.Register(ModuleName, 2011, "valid bet UID should be provided")
	ErrInvalidOddsUID             = sdkerrors.Register(ModuleName, 2013, "valid odds UID should be provided")
	ErrInvalidOddsValue           = sdkerrors.Register(ModuleName, 2014, "valid odds value should be provided (odds value more than 1.0)")
	ErrInvalidSportEventUID       = sdkerrors.Register(ModuleName, 2015, "valid sport event UID should be provided")
	ErrInvalidTicket              = sdkerrors.Register(ModuleName, 2016, "valid ticket should be provided")
	ErrInvalidAmount              = sdkerrors.Register(ModuleName, 2017, "valid amount should be provided")
	ErrEmptyUidsList              = sdkerrors.Register(ModuleName, 2018, "empty UIDs list")
	ErrTooManyUids                = sdkerrors.Register(ModuleName, 2019, "too many UIDs")
	ErrInJSONMarshal              = sdkerrors.Register(ModuleName, 2020, "internal error in Json marshaling")
	ErrInMarshaling               = sdkerrors.Register(ModuleName, 2021, "internal error in marshaling")
	ErrNoMatchingBet              = sdkerrors.Register(ModuleName, 2022, "no matching bet")
	ErrResultNotDeclared          = sdkerrors.Register(ModuleName, 2023, "sport event result is not declared")
	ErrDuplicateUID               = sdkerrors.Register(ModuleName, 2024, "UID is already set")
	ErrTooManyBets                = sdkerrors.Register(ModuleName, 2025, "too many bets")
	ErrEmptyBetListRequest        = sdkerrors.Register(ModuleName, 2026, "no bet in the request")
	ErrInSRPlacementProcessing    = sdkerrors.Register(ModuleName, 2027, "internal error in processing bet placement in SR")
	ErrSportEventStatusNotPending = sdkerrors.Register(ModuleName, 2028, "can not place bet on this sport event any more")
	ErrEndTSIsPassed              = sdkerrors.Register(ModuleName, 2029, "the sport event is expired")
	ErrOddsUIDNotExist            = sdkerrors.Register(ModuleName, 2030, "sport event does not have this odds UID")
	ErrSportEventIsAborted        = sdkerrors.Register(ModuleName, 2031, "sport event is aborted")
	ErrInvalidCreatorAddr         = sdkerrors.Register(ModuleName, 2032, "can not cretae an AccAddress from creator")
	ErrInSRRefund                 = sdkerrors.Register(ModuleName, 2033, "internal error in refunding user in SR")
	ErrInactiveSportEvent         = sdkerrors.Register(ModuleName, 2034, "the sport event is not active")
	ErrBetAmountIsLow             = sdkerrors.Register(ModuleName, 2035, "bet amount is lower than the minimum allowed")
	ErrInAddAmountToSportEvent    = sdkerrors.Register(ModuleName, 2036, "internal error in adding bet amount to sport event")
	ErrInConvertingOddsToDec      = sdkerrors.Register(ModuleName, 2037, "internal error in converting odds value from string to sdk.Dec")
	ErrInSubAmountFromSportEvent  = sdkerrors.Register(ModuleName, 2038, "internal error in adding bet amount to sport event")
	ErrVigIsOutOfRange            = sdkerrors.Register(ModuleName, 2039, "vig is out of valid range")
	ErrEventMaxLossExceeded       = sdkerrors.Register(ModuleName, 2040, "event loss is in max state")
	ErrOddsMaxLossExceeded        = sdkerrors.Register(ModuleName, 2041, "odds loss is in max state")
	ErrOddsDataNotFound           = sdkerrors.Register(ModuleName, 2042, "odds not exist in ticket payload")
	ErrTextInvalidCreator         = "invalid creator address (%s)"
)
