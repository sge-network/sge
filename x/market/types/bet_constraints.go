package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// validate validates the bet constraints
func (bc *MarketBetConstraints) validate(params *Params) error {
	if bc.BetFee.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market bet fee can not be negative")
	}

	if bc.BetFee.LT(params.MarketMinBetFee) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market bet fee is out of threshold limit")
	}

	if bc.MinAmount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market min amount can not be negative")
	}

	if bc.MinAmount.LT(params.MarketMinBetAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "market min bet amount is less than threshold")
	}

	return nil
}
