package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AffiliateReward is the type for affiliation rewards calculations
type AffiliateReward struct{}

// NewAffiliateReward create new object of affiliation reward calculator type.
func NewAffiliateReward() AffiliateReward { return AffiliateReward{} }

// VaidateCampaign validates campaign definitions.
func (rfr AffiliateReward) VaidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_SIGNUP {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "affiliate rewards can only have single definition")
	}
	if campaign.RewardAmount.MainAccountAmount.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrInvalidGranteeType, "affiliate rewards can be defined for subaccount only")
	}
	if campaign.RewardAmount.SubaccountAmount.LTE(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "affiliate rewards for subaccount should be positive")
	}

	// TODO: validate duplicate signup affiliate reward
	return nil
}

// Calculate parses ticket payload and returns the distribution list of referral reward.
func (rfr AffiliateReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket string,
) (Allocation, error) {
	var payload GrantSignupRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return Allocation{}, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	_, found := keepers.SubAccountKeeper.GetSubAccountOwner(ctx, sdk.MustAccAddressFromBech32(payload.Common.Receiver))
	if !found {
		return Allocation{}, sdkerrors.Wrapf(ErrReceiverAddrNotSubAcc, "%s", &payload.Common.Receiver)
	}

	// TODO: validate reward grant for referral

	return Allocation{
		SubAcc: NewReceiver(
			payload.Common.Receiver,
			campaign.RewardAmount.SubaccountAmount,
			0,
		),
	}, nil
}
