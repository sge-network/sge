package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewEventBetConstraints creates new bet constraint pointer
func NewEventBetConstraints(minAmount, betFee sdk.Int) *EventBetConstraints {
	return &EventBetConstraints{
		MinAmount: minAmount,
		BetFee:    betFee,
	}
}

// validate validates the bet constraints
func (bc *EventBetConstraints) validate(params *Params) error {
	if bc.BetFee.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event bet fee can not be negative")
	}

	if bc.BetFee.LT(params.EventMinBetFee) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event bet fee is out of threshold limit")
	}

	if bc.MinAmount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min amount can not be negative")
	}

	if bc.MinAmount.LT(params.EventMinBetAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min bet amount is less than threshold")
	}

	return nil
}
