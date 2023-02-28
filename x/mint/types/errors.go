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
	ErrTextInvalidParamType                   = "invalid parameter type: %T"
	ErrTextBlocksPerYearMustBePositive        = "blocks per year must be positive: %d"
	ErrTextExcludeAmountMustBePositive        = "exclude amount must be positive: %s"
	ErrTextPhasesShouldHaveValue              = "phases should have value: %v"
	ErrTextMintParamInflationShouldBePositive = "mint parameter Inflation should be positive, is %s"
	ErrTextYearCoefficientMustBePositive      = "year coefficient should be non-zero and positive value"
	ErrTextEndPhaseParamNotAllowed            = "adding phase with equal values with end phase is not allowed"
	ErrTextNilMinter                          = "stored minter should not be nil"
)
