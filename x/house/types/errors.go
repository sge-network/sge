package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/house module sentinel errors
var (
	ErrEmptyDepositorAddr        = sdkerrors.Register(ModuleName, 5001, "empty depositor address")
	ErrInvalidMarketUID          = sdkerrors.Register(ModuleName, 5002, "invalid market uid")
	ErrInvalidWithdrawMode       = sdkerrors.Register(ModuleName, 5003, "invalid withdrawal mode")
	ErrInvalidIndex              = sdkerrors.Register(ModuleName, 5004, "invalid participant index")
	ErrInvalidAmount             = sdkerrors.Register(ModuleName, 5005, "invalid amount")
	ErrDepositTooSmall           = sdkerrors.Register(ModuleName, 5006, "deposit amount is less than minimum acceptable deposit")
	ErrOBDepositProcessing       = sdkerrors.Register(ModuleName, 5007, "internal error in processing deposit in orderbook")
	ErrDepositSetting            = sdkerrors.Register(ModuleName, 5008, "internal error in setting deposit")
	ErrInvalidparticipationIndex = sdkerrors.Register(ModuleName, 5009, "invalid participant index")
	ErrInvalidMode               = sdkerrors.Register(ModuleName, 5010, "invalid withdrawal mode")
	ErrDepositNotFound           = sdkerrors.Register(ModuleName, 5011, "deposit not found")
	ErrWithdrawalTooLarge        = sdkerrors.Register(ModuleName, 5012, "withdrawal is more than unused amount")
	ErrOBLiquidateProcessing     = sdkerrors.Register(ModuleName, 5013, "internal error in processing liquidation in orderbook")
	ErrWrongWithdrawCreator      = sdkerrors.Register(ModuleName, 5014, "withdrawal is only allowed from the depositor account")
	ErrInTicketVerification      = sdkerrors.Register(ModuleName, 5015, "error in ticket verification process")
	ErrInTicketPayloadValidation = sdkerrors.Register(ModuleName, 5016, "error in ticket payload validation")
	ErrUserKycFailed             = sdkerrors.Register(ModuleName, 5017, "the account failed the KYC Validation")
	ErrAuthorizationNotFound     = sdkerrors.Register(ModuleName, 5018, "no authorization found")
	ErrAuthorizationNotAccepted  = sdkerrors.Register(ModuleName, 5019, "authorization not accepted")
)
