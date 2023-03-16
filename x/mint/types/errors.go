package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/mint module sentinel errors
var (
	ErrMintDenomIsBlank = sdkerrors.Register(ModuleName, 3001, "mint denom cannot be blank")
)

const (
	ErrTextInvalidParamType                   = "invalid parameter type"
	ErrTextBlocksPerYearMustBePositive        = "blocks per year must be positive"
	ErrTextExcludeAmountMustBePositive        = "exclude amount must be positive"
	ErrTextPhasesShouldHaveValue              = "phases should have value"
	ErrTextMintParamInflationShouldBePositive = "mint parameter Inflation should be positive"
	ErrTextYearCoefficientMustBePositive      = "year coefficient should be non-zero and positive value"
	ErrTextEndPhaseParamNotAllowed            = "adding phase with equal values with end phase is not allowed"
	ErrTextNilMinter                          = "stored minter should not be nil"
)
