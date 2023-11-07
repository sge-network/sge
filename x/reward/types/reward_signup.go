package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SignUpReward is the type for signup rewards calculations
type SignUpReward struct{}

// NewSignUpReward create new object of signup reward calculator type.
func NewSignUpReward() SignUpReward { return SignUpReward{} }

// VaidateCampaign validates campaign definitions.
func (sur SignUpReward) VaidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_SIGNUP {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "signup rewards can only have single definition")
	}
	if campaign.RewardAmount.MainAccountAmount.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrInvalidGranteeType, "signup rewards can be defined for subaccount only")
	}
	if campaign.RewardAmount.SubaccountAmount.LTE(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "signup rewards for subaccount should be positive")
	}

	// TODO: validate duplicate signup reward
	return nil
}

// Calculate parses ticket payload and returns the distribution list of signup reward.
func (sur SignUpReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket string,
) (Receiver, RewardPayloadCommon, bool, string, error) {
	var payload GrantSignupRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return Receiver{}, payload.Common, false, "", sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	_, found := keepers.SubAccountKeeper.GetSubAccountOwner(ctx, sdk.MustAccAddressFromBech32(payload.Common.Receiver))
	if !found {
		return Receiver{}, payload.Common, false, "", sdkerrors.Wrapf(ErrReceiverAddrNotSubAcc, "%s", &payload.Common.Receiver)
	}

	// TODO: validate reward grant for sighnup

	return NewReceiver(
		payload.Common.Receiver,
		campaign.RewardAmount.SubaccountAmount,
		0,
	), payload.Common, true, payload.Common.Receiver, nil
}
