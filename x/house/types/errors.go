package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/house module sentinel errors
var (
	ErrEmptyDepositorAddr           = sdkerrors.Register(ModuleName, 5001, "empty depositor address")
	ErrInvalidSportEventUID         = sdkerrors.Register(ModuleName, 5002, "invalid sport event uid")
	ErrInvalidWithdrawMode          = sdkerrors.Register(ModuleName, 5003, "invalid withdrawal mode")
	ErrInvalidIndex                 = sdkerrors.Register(ModuleName, 5004, "invalid participant index")
	ErrInvalidAmount                = sdkerrors.Register(ModuleName, 5005, "valid amount should be provided")
	ErrDepositTooSmall              = sdkerrors.Register(ModuleName, 5006, "deposit is not greater than minimum deposit amount")
	ErrOrderBookDepositProcessing   = sdkerrors.Register(ModuleName, 5007, "internal error in processing deposit in OB")
	ErrDepositSetting               = sdkerrors.Register(ModuleName, 5008, "internal error in setting deposit")
	ErrInvalidparticipationIndex    = sdkerrors.Register(ModuleName, 5009, "invalid participant index")
	ErrInvalidMode                  = sdkerrors.Register(ModuleName, 5010, "invalid withdrawal mode")
	ErrDepositNotFound              = sdkerrors.Register(ModuleName, 5011, "deposit not found")
	ErrWithdrawalTooLarge           = sdkerrors.Register(ModuleName, 5012, "withdrawal is more than unused amount")
	ErrOrderBookLiquidateProcessing = sdkerrors.Register(ModuleName, 5013, "internal error in processing liquidation in OB")
	ErrWrongWithdrawCreator         = sdkerrors.Register(ModuleName, 5014, "withdrawal is only allowed from the depositor account")
)

const (
	ErrTextInvalidDepositor = "invalid depositor address (%s)"
)
