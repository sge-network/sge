package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/market module sentinel errors
var (
	ErrMarketCanNotBeAltered           = sdkerrors.Register(ModuleName, 1001, "market cannot be altered if it is not active or inactive")
	ErrMarketAlreadyExist              = sdkerrors.Register(ModuleName, 1002, "market already exist")
	ErrMarketNotFound                  = sdkerrors.Register(ModuleName, 1003, "market not found")
	ErrMarketResolutionNotAllowed      = sdkerrors.Register(ModuleName, 1004, "market resolution is allowed for active or inactive status")
	ErrInvalidWinnerOdds               = sdkerrors.Register(ModuleName, 1005, "the provided winner odds does not exist in the market odds")
	ErrInTicketVerification            = sdkerrors.Register(ModuleName, 1006, "error in ticket verification process")
	ErrInTicketPayloadValidation       = sdkerrors.Register(ModuleName, 1007, "error in ticket payload validation")
	ErrResolutionTimeLessThenStartTime = sdkerrors.Register(ModuleName, 1008, "resolution time cannot be less than market start time")
	ErrInOrderBookInitiation           = sdkerrors.Register(ModuleName, 1009, "error in order book initiation")
)
