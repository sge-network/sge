package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &CreateCampaignAuthorization{}
var _ authz.Authorization = &UpdateCampaignAuthorization{}

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
func (a CreateCampaignAuthorization) Accept(_ sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
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
	if a.SpendLimit.LTE(sdk.ZeroInt()) {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be less than or equal to zero")
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
func (a UpdateCampaignAuthorization) Accept(_ sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgUpdateCampaign, ok := msg.(*MsgUpdateCampaign)
	if !ok {
		return authz.AcceptResponse{}, sdkerrtypes.ErrInvalidType.Wrap("type mismatch")
	}

	limitLeft := a.SpendLimit
	if msgUpdateCampaign.TopupFunds.GT(sdkmath.ZeroInt()) {
		limitLeft := limitLeft.Sub(msgUpdateCampaign.TopupFunds)
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
	if a.SpendLimit.LTE(sdk.ZeroInt()) {
		return sdkerrtypes.ErrInvalidCoins.Wrap("spend limit cannot be less than or equal to zero")
	}

	return nil
}
