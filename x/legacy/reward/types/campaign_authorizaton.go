package types

import (
	context "context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var (
	_ authz.Authorization = &CreateCampaignAuthorization{}
	_ authz.Authorization = &UpdateCampaignAuthorization{}
	_ authz.Authorization = &WithdrawCampaignAuthorization{}
)

// NewCreateCampaignAuthorization creates a new CreateCampaignAuthorization object.
func NewCreateCampaignAuthorization(spendLimit sdkmath.Int) *CreateCampaignAuthorization {
	return &CreateCampaignAuthorization{
		SpendLimit: spendLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (CreateCampaignAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgCreateCampaign{})
}

// Accept implements Authorization.Accept.
func (a CreateCampaignAuthorization) Accept(_ context.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgCreateCampaign, ok := msg.(*MsgCreateCampaign)
	if !ok {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInvalidType.Wrap("type mismatch")
	}

	limitLeft := a.SpendLimit.Sub(msgCreateCampaign.TotalFunds)
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
		Updated: &CreateCampaignAuthorization{SpendLimit: limitLeft},
	}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a CreateCampaignAuthorization) ValidateBasic() error {
	if a.SpendLimit.IsNil() {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be nil")
	}
	if a.SpendLimit.LTE(sdkmath.ZeroInt()) {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be less than or equal to zero")
	}

	if a.SpendLimit.LT(minCampaignFunds) {
		return sdkerrtypes.ErrInvalidCoins.Wrapf("spend limit cannot be less than %s", minCampaignFunds)
	}

	return nil
}

// NewUpdateCampaignAuthorization creates a new UpdateCampaignAuthorization object.
func NewUpdateCampaignAuthorization(spendLimit sdkmath.Int) *UpdateCampaignAuthorization {
	return &UpdateCampaignAuthorization{
		SpendLimit: spendLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (UpdateCampaignAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgUpdateCampaign{})
}

// Accept implements Authorization.Accept.
func (a UpdateCampaignAuthorization) Accept(_ context.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgUpdateCampaign, ok := msg.(*MsgUpdateCampaign)
	if !ok {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInvalidType.Wrap("type mismatch")
	}

	limitLeft := a.SpendLimit
	if msgUpdateCampaign.TopupFunds.GT(sdkmath.ZeroInt()) {
		limitLeft = limitLeft.Sub(msgUpdateCampaign.TopupFunds)
		if limitLeft.IsNegative() {
			return authz.AcceptResponse{}, sdkerrtypes.ErrInsufficientFunds.Wrapf(
				"requested amount is more than spend limit",
			)
		}
	}

	if limitLeft.IsZero() {
		return authz.AcceptResponse{Accept: true, Delete: true}, nil
	}

	return authz.AcceptResponse{
		Accept:  true,
		Delete:  false,
		Updated: &UpdateCampaignAuthorization{SpendLimit: limitLeft},
	}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a UpdateCampaignAuthorization) ValidateBasic() error {
	if a.SpendLimit.IsNil() {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be nil")
	}
	if a.SpendLimit.LTE(sdkmath.ZeroInt()) {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be less than or equal to zero")
	}

	if a.SpendLimit.LT(minCampaignFunds) {
		return sdkerrtypes.ErrInvalidCoins.Wrapf("spend limit cannot be less than %s", minCampaignFunds)
	}

	return nil
}

var _ authz.Authorization = &WithdrawCampaignAuthorization{}

// NewWithdrawAuthorization creates a new WithdrawAuthorization object.
func NewWithdrawAuthorization(withdrawLimit sdkmath.Int) *WithdrawCampaignAuthorization {
	return &WithdrawCampaignAuthorization{
		WithdrawLimit: withdrawLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (WithdrawCampaignAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgWithdrawFunds{})
}

// Accept implements Authorization.Accept.
func (a WithdrawCampaignAuthorization) Accept(_ context.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	mWithdraw, ok := msg.(*MsgWithdrawFunds)
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
		Updated: &WithdrawCampaignAuthorization{WithdrawLimit: limitLeft},
	}, nil
}

// ValidateBasic implements Authorization.ValidateBasic.
func (a WithdrawCampaignAuthorization) ValidateBasic() error {
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
