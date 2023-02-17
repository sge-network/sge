package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/house module sentinel errors
var (
	ErrEmptyDepositorAddr           = sdkerrors.Register(ModuleName, 1001, "empty depositor address")
	ErrInvalidSportEventUID         = sdkerrors.Register(ModuleName, 1002, "invalid sport event uid")
	ErrInvalidDenom                 = sdkerrors.Register(ModuleName, 1003, "invalid coin denomination")
	ErrDepositTooSmall              = sdkerrors.Register(ModuleName, 1004, "deposit is not greater than minimum deposit amount")
	ErrOrderBookDepositProcessing   = sdkerrors.Register(ModuleName, 1005, "internal error in processing deposit in OB")
	ErrDepositSetting               = sdkerrors.Register(ModuleName, 1006, "internal error in setting deposit")
	ErrInvalidParticipantID         = sdkerrors.Register(ModuleName, 1007, "invalid participant id")
	ErrInvalidMode                  = sdkerrors.Register(ModuleName, 1008, "invalid withdrawal mode")
	ErrDepositNotFound              = sdkerrors.Register(ModuleName, 1009, "deposit not found")
	ErrWithdrawalTooLarge           = sdkerrors.Register(ModuleName, 1010, "withdrawal is more than unused amount")
	ErrOrderBookLiquidateProcessing = sdkerrors.Register(ModuleName, 1011, "internal error in processing liquidation in OB")
)

const (
	ErrTextInvalidDepositor = "invalid depositor address (%s)"
)
