package types

import (
	context "context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &WithdrawAuthorization{}

// NewWithdrawAuthorization creates a new WithdrawAuthorization object.
func NewWithdrawAuthorization(withdrawLimit sdkmath.Int) *WithdrawAuthorization {
	return &WithdrawAuthorization{
		WithdrawLimit: withdrawLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (WithdrawAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgWithdraw{})
}

// Accept implements Authorization.Accept.
func (a WithdrawAuthorization) Accept(_ context.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	mWithdraw, ok := msg.(*MsgWithdraw)
	if !ok {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInvalidType.Wrap("type mismatch")
	}

	limitLeft := a.WithdrawLimit.Sub(mWithdraw.Amount)
	if limitLeft.IsNegative() {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInsufficientFunds.Wrapf(
			"requested amount is more than withdraw limit",
		)
	}
	if limitLeft.IsZero() {
		return authz.AcceptResponse{Accept: true, Delete: true}, nil
	}

	return authz.AcceptResponse{
		Accept:  true,
		Delete:  false,
		Updated: &WithdrawAuthorization{WithdrawLimit: limitLeft},
	}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a WithdrawAuthorization) ValidateBasic() error {
	if a.WithdrawLimit.IsNil() {
		return sdkerrtypes.ErrInvalidCoins.Wrap("withdraw limit cannot be nil")
	}
	if a.WithdrawLimit.LTE(sdkmath.ZeroInt()) {
		return sdkerrtypes.ErrInvalidCoins.Wrap("withdraw limit cannot be less than or equal to zero")
	}
	if a.WithdrawLimit.GT(maxWithdrawGrant) {
		return sdkerrtypes.ErrInvalidCoins.Wrapf(
			"withdraw limit cannot be grated than %s",
			maxWithdrawGrant,
		)
	}

	return nil
}
