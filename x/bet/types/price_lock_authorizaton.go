package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &PriceLockAuthorization{}

// NewPriceLockAuthorization creates a new PriceLockAuthorization object.
func NewPriceLockAuthorization(spendLimit sdkmath.Int) *PriceLockAuthorization {
	return &PriceLockAuthorization{
		SpendLimit: spendLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (PriceLockAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgPriceLockPoolTopUp{})
}

// Accept implements Authorization.Accept.
func (a PriceLockAuthorization) Accept(_ sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	mPriceLockFund, ok := msg.(*MsgPriceLockPoolTopUp)
	if !ok {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInvalidType.Wrap("type mismatch")
	}

	limitLeft := a.SpendLimit.Sub(mPriceLockFund.Amount)
	if limitLeft.IsNegative() {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInsufficientFunds.Wrapf(
			"requested amount is more than spend limit",
		)
	}
	if limitLeft.IsZero() {
		return authz.AcceptResponse{Accept: true, Delete: true}, nil
	}

	return authz.AcceptResponse{
		Accept:  true,
		Delete:  false,
		Updated: &PriceLockAuthorization{SpendLimit: limitLeft},
	}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a PriceLockAuthorization) ValidateBasic() error {
	if a.SpendLimit.IsNil() {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be nil")
	}
	if a.SpendLimit.LTE(sdk.ZeroInt()) {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be less than or equal to zero")
	}
	if a.SpendLimit.LT(minPriceLockFund) {
		return sdkerrtypes.ErrInvalidCoins.Wrapf("spend limit cannot be less than %s", minPriceLockFund)
	}

	return nil
}
