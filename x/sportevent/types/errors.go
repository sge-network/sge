package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sportevent module sentinel errors
var (
	ErrInMarshaling               = sdkerrors.Register(ModuleName, 1001, "internal error in marshaling")
	ErrInUnmarshaling             = sdkerrors.Register(ModuleName, 1002, "internal error in unmarshaling")
	ErrNoMatchingSportEvent       = sdkerrors.Register(ModuleName, 1003, "sport-event not found")
	ErrEmptySportEventUIDs        = sdkerrors.Register(ModuleName, 1004, "sport-event UID(s) must be provided")
	ErrCanNotBeAltered            = sdkerrors.Register(ModuleName, 1006, "event cannot be altered if it is active or inactive")
	ErrEventAlreadyExist          = sdkerrors.Register(ModuleName, 1007, "event already exist")
	ErrEventNotFound              = sdkerrors.Register(ModuleName, 1008, "event not found")
	ErrEventIsNotActiveOrInactive = sdkerrors.Register(ModuleName, 1009, "event cannot be resolved as status is not active or inactive")
	ErrInvalidWinnerOdds          = sdkerrors.Register(ModuleName, 1010, "the provided winner odds does not exist in the event odds")
	ErrInVerification             = sdkerrors.Register(ModuleName, 1012, "error in verification process")
	ErrResolutionTimeLessTnStart  = sdkerrors.Register(ModuleName, 1013, "resolution time cannot be less than event start time")
	ErrDuplicateEventsExist       = sdkerrors.Register(ModuleName, 1014, "duplicate events provided")
	ErrOddsStatsNotFound          = sdkerrors.Register(ModuleName, 1015, "odds stats not found")
	ErrInOrderBookInitiation      = sdkerrors.Register(ModuleName, 1016, "error in order book initiation")
)
