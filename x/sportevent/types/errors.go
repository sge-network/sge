package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sportevent module sentinel errors
var (
	ErrInMarshaling              = sdkerrors.Register(ModuleName, 3001, "Internal error in marshaling")
	ErrInUnmarshaling            = sdkerrors.Register(ModuleName, 3002, "Internal error in unmarshaling")
	ErrNoMatchingSportEvent      = sdkerrors.Register(ModuleName, 3003, "No any matching sport event")
	ErrEmptySportEventUIDs       = sdkerrors.Register(ModuleName, 3004, "Sport event UID(s) should be provided")
	ErrCanNotBeAltered           = sdkerrors.Register(ModuleName, 3006, "event cannot be altered after resolution")
	ErrEventAlreadyExist         = sdkerrors.Register(ModuleName, 3007, "event already exist")
	ErrEventNotFound             = sdkerrors.Register(ModuleName, 3008, "event not found")
	ErrEventIsNotPending         = sdkerrors.Register(ModuleName, 3009, "event cannot be resolved as status is not pending")
	ErrInvalidWinnerOdd          = sdkerrors.Register(ModuleName, 3010, "the provided winner odd not exist in the event odds")
	ErrMaxBetCapExceeded         = sdkerrors.Register(ModuleName, 3011, "provided amount exceeds the event bet capacity")
	ErrInVerification            = sdkerrors.Register(ModuleName, 3012, "error in verification process")
	ErrResolutionTimeLessTnStart = sdkerrors.Register(ModuleName, 3013, "resolution time cannot be less than start")
	ErrDuplicateEventsExist      = sdkerrors.Register(ModuleName, 30014, "duplicate events provided")
	ErrOddsStatsNotFound         = sdkerrors.Register(ModuleName, 30015, "odds stats not found")
	ErrEventMaxLossExceeded      = sdkerrors.Register(ModuleName, 30016, "max loss of the event exceeded")
)
