package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &WithdrawAuthorization{}

// NewWithdrawAuthorization creates a new WithdrawAuthorization object.
func NewWithdrawAuthorization(withdrawLimit sdk.Int) *WithdrawAuthorization {
	return &WithdrawAuthorization{
		WithdrawLimit: withdrawLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (a WithdrawAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgDeposit{})
}

// Accept implements Authorization.Accept.
func (a WithdrawAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	mDeposit, ok := msg.(*MsgDeposit)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}

	if a.WithdrawLimit.LT(mDeposit.Amount) {
		return authz.AcceptResponse{}, sdkerrors.ErrInsufficientFunds.Wrapf("requested amount is more than withdraw limit")
	}

	return authz.AcceptResponse{Accept: true, Delete: false, Updated: &WithdrawAuthorization{WithdrawLimit: mDeposit.Amount}}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a WithdrawAuthorization) ValidateBasic() error {
	if a.WithdrawLimit.IsNil() {
		return sdkerrors.ErrInvalidCoins.Wrap("withdraw limit cannot be nil")
	}
	if a.WithdrawLimit.LTE(sdk.ZeroInt()) {
		return sdkerrors.ErrInvalidCoins.Wrap("withdraw limit cannot be less than or equal to zero")
	}

	return nil
}
