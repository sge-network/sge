package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/house module sentinel errors
var (
	ErrInTicketVerification      = sdkerrors.Register(ModuleName, 5001, "error in ticket verification process")
	ErrInTicketPayloadValidation = sdkerrors.Register(ModuleName, 5002, "error in ticket payload validation")
	ErrInvalidMarketUID          = sdkerrors.Register(ModuleName, 5003, "invalid market uid")
	ErrInvalidWithdrawMode       = sdkerrors.Register(ModuleName, 5004, "invalid withdrawal mode")
	ErrInvalidIndex              = sdkerrors.Register(ModuleName, 5005, "invalid participant index")
	ErrInvalidAmount             = sdkerrors.Register(ModuleName, 5006, "invalid amount")
	ErrDepositTooSmall           = sdkerrors.Register(ModuleName, 5007, "deposit amount is less than minimum acceptable deposit")
	ErrOBDepositProcessing       = sdkerrors.Register(ModuleName, 5008, "internal error in processing deposit in orderbook")
	ErrInvalidMode               = sdkerrors.Register(ModuleName, 5009, "invalid withdrawal mode")
	ErrDepositNotFound           = sdkerrors.Register(ModuleName, 5010, "deposit not found")
	ErrOBLiquidateProcessing     = sdkerrors.Register(ModuleName, 5011, "internal error in processing liquidation in orderbook")
	ErrUserKycFailed             = sdkerrors.Register(ModuleName, 5012, "the account failed the KYC Validation")
	ErrAuthorizationNotFound     = sdkerrors.Register(ModuleName, 5013, "no authorization found")
	ErrAuthorizationNotAccepted  = sdkerrors.Register(ModuleName, 5014, "authorization not accepted")
	ErrAuthorizationNotAllowed   = sdkerrors.Register(ModuleName, 5015, "authorization not allowed")
)
