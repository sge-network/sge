package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/house module sentinel errors
// nolint
var (
	ErrEmptyDelegatorAddr         = sdkerrors.Register(ModuleName, 1001, "empty depositor address")
	ErrInvalidSportEventUid       = sdkerrors.Register(ModuleName, 1002, "invalid sport event uid")
	ErrInvalidDenom               = sdkerrors.Register(ModuleName, 1003, "invalid coin denomination")
	ErrDepositTooSmall            = sdkerrors.Register(ModuleName, 1004, "deposit is not greater than minimum deposit amount")
	ErrOrderBookDepositProcessing = sdkerrors.Register(ModuleName, 1005, "internal error in processing deposit in OB")
	ErrDepositSetting             = sdkerrors.Register(ModuleName, 1006, "internal error in setting deposit")
)

const (
	ErrTextInvalidDepositor = "invalid depositor address (%s)"
)
