package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// validate validates the bet constraints
func (bc *MarketBetConstraints) validate(params *Params) error {
	if bc.BetFee.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market bet fee can not be negative")
	}

	if bc.BetFee.LT(params.MinBetFee) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market bet fee is lower than minimum threshold limit")
	}

	if bc.BetFee.GT(params.MaxBetFee) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market bet fee is higher than the threshold limit")
	}

	if bc.MinAmount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market min amount can not be negative")
	}

	if bc.MinAmount.LT(params.MinBetAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market min bet amount is less than threshold")
	}

	return nil
}
