package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &DepositAuthorization{}

// NewDepositAuthorization creates a new DepositAuthorization object.
func NewDepositAuthorization(spendLimit sdk.Int) *DepositAuthorization {
	return &DepositAuthorization{
		SpendLimit: spendLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (a DepositAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgDeposit{})
}

// Accept implements Authorization.Accept.
func (a DepositAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	mDeposit, ok := msg.(*MsgDeposit)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}

	if a.SpendLimit.LT(mDeposit.Amount) {
		return authz.AcceptResponse{}, sdkerrors.ErrInsufficientFunds.Wrapf("requested amount is more than spend limit")
	}

	return authz.AcceptResponse{Accept: true, Delete: false, Updated: &DepositAuthorization{SpendLimit: mDeposit.Amount}}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a DepositAuthorization) ValidateBasic() error {
	if a.SpendLimit.IsNil() {
		return sdkerrors.ErrInvalidCoins.Wrap("spend limit cannot be nil")
	}
	if a.SpendLimit.LTE(sdk.ZeroInt()) {
		return sdkerrors.ErrInvalidCoins.Wrap("spend limit cannot be less than or equal to zero")
	}

	return nil
}
