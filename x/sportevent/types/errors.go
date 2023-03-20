package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sportevent module sentinel errors
var (
	ErrNoMatchingSportEvent       = sdkerrors.Register(ModuleName, 1001, "sport-event not found")
	ErrCanNotBeAltered            = sdkerrors.Register(ModuleName, 1002, "event cannot be altered if it is active or inactive")
	ErrEventAlreadyExist          = sdkerrors.Register(ModuleName, 1003, "event already exist")
	ErrEventNotFound              = sdkerrors.Register(ModuleName, 1004, "event not found")
	ErrEventIsNotActiveOrInactive = sdkerrors.Register(ModuleName, 1005, "event cannot be resolved as status is not active or inactive")
	ErrInvalidWinnerOdds          = sdkerrors.Register(ModuleName, 1006, "the provided winner odds does not exist in the event odds")
	ErrInVerification             = sdkerrors.Register(ModuleName, 1007, "error in verification process")
	ErrResolutionTimeLessTnStart  = sdkerrors.Register(ModuleName, 1008, "resolution time cannot be less than event start time")
	ErrInOrderBookInitiation      = sdkerrors.Register(ModuleName, 1009, "error in order book initiation")
)
