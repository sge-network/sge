package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReferralReward is the type for referral rewards calculations
type ReferralReward struct{}

// NewReferralReward create new object of referral reward calculator type.
func NewReferralReward() ReferralReward { return ReferralReward{} }

// VaidateCampaign validates campaign definitions.
func (rfr ReferralReward) VaidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_SIGNUP {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "referral rewards can only have single definition")
	}
	if campaign.RewardAmount.MainAccountAmount.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrInvalidGranteeType, "referral rewards can be defined for subaccount only")
	}
	if campaign.RewardAmount.SubaccountAmount.LTE(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "referral rewards for subaccount should be positive")
	}

	// TODO: validate duplicate signup referral reward
	return nil
}

// Calculate parses ticket payload and returns the distribution list of referral reward.
func (rfr ReferralReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
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

	// TODO: validate reward grant for referral

	return NewReceiver(
		payload.Common.Receiver,
		campaign.RewardAmount.SubaccountAmount,
		0,
	), payload.Common, true, "", nil
}
