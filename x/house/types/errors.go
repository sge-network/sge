package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/house module sentinel errors
var (
	ErrEmptyDepositorAddr           = sdkerrors.Register(ModuleName, 5001, "empty depositor address")
	ErrInvalidSportEventUID         = sdkerrors.Register(ModuleName, 5002, "invalid sport event uid")
	ErrInvalidAmount                = sdkerrors.Register(ModuleName, 5003, "valid amount should be provided")
	ErrDepositTooSmall              = sdkerrors.Register(ModuleName, 5004, "deposit is not greater than minimum deposit amount")
	ErrOrderBookDepositProcessing   = sdkerrors.Register(ModuleName, 5005, "internal error in processing deposit in OB")
	ErrDepositSetting               = sdkerrors.Register(ModuleName, 5006, "internal error in setting deposit")
	ErrInvalidparticipationIndex    = sdkerrors.Register(ModuleName, 5007, "invalid participant index")
	ErrInvalidMode                  = sdkerrors.Register(ModuleName, 5008, "invalid withdrawal mode")
	ErrDepositNotFound              = sdkerrors.Register(ModuleName, 5009, "deposit not found")
	ErrWithdrawalTooLarge           = sdkerrors.Register(ModuleName, 5010, "withdrawal is more than unused amount")
	ErrOrderBookLiquidateProcessing = sdkerrors.Register(ModuleName, 5011, "internal error in processing liquidation in OB")
)

const (
	ErrTextInvalidDepositor = "invalid depositor address (%s)"
)
