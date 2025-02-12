package types

import (
	sdkmath "cosmossdk.io/math"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates fields of the given ticketData
func (payload *SubAccWagerTicketPayload) Validate(betAmount sdkmath.Int) error {
	if payload.MainaccDeductAmount.IsNil() || payload.SubaccDeductAmount.IsNil() {
		return sdkerrtypes.ErrInvalidRequest.Wrap("main account and subaccount deduction should be set")
	}

	if !payload.MainaccDeductAmount.Add(payload.SubaccDeductAmount).Equal(betAmount) {
		return sdkerrtypes.ErrInvalidRequest.Wrap("sum of main and sub account deduction should be equal to bet amount")
	}

	return payload.Msg.ValidateBasic()
}
