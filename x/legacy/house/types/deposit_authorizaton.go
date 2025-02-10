package types

import (
	context "context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &DepositAuthorization{}

// NewDepositAuthorization creates a new DepositAuthorization object.
func NewDepositAuthorization(spendLimit sdkmath.Int) *DepositAuthorization {
	return &DepositAuthorization{
		SpendLimit: spendLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (DepositAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgDeposit{})
}

// Accept implements Authorization.Accept.
func (a DepositAuthorization) Accept(_ context.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	mDeposit, ok := msg.(*MsgDeposit)
	if !ok {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInvalidType.Wrap("type mismatch")
	}

	limitLeft := a.SpendLimit.Sub(mDeposit.Amount)
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
		Updated: &DepositAuthorization{SpendLimit: limitLeft},
	}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a DepositAuthorization) ValidateBasic() error {
	if a.SpendLimit.IsNil() {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be nil")
	}
	if a.SpendLimit.LTE(sdkmath.ZeroInt()) {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be less than or equal to zero")
	}
	if a.SpendLimit.LT(minDepositGrant) {
		return sdkerrtypes.ErrInvalidCoins.Wrapf("spend limit cannot be less than %s", minDepositGrant)
	}

	return nil
}
