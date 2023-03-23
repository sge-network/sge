package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/market module sentinel errors
var (
	ErrNoMatchingMarket                = sdkerrors.Register(ModuleName, 1001, "market not found")
	ErrMarketCanNotBeAltered           = sdkerrors.Register(ModuleName, 1002, "market cannot be altered if it is active or inactive")
	ErrMarketAlreadyExist              = sdkerrors.Register(ModuleName, 1003, "market already exist")
	ErrMarketNotFound                  = sdkerrors.Register(ModuleName, 1004, "market not found")
	ErrMarketIsNotActiveOrInactive     = sdkerrors.Register(ModuleName, 1005, "market cannot be resolved as status is not active or inactive")
	ErrInvalidWinnerOdds               = sdkerrors.Register(ModuleName, 1006, "the provided winner odds does not exist in the market odds")
	ErrInVerification                  = sdkerrors.Register(ModuleName, 1007, "error in verification process")
	ErrResolutionTimeLessThenStartTime = sdkerrors.Register(ModuleName, 1008, "resolution time cannot be less than market start time")
	ErrInOrderBookInitiation           = sdkerrors.Register(ModuleName, 1009, "error in order book initiation")
)
