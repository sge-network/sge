package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// SignUpReferrerReward is the type for signup referrer rewards calculations
type SignUpReferrerReward struct{}

// NewSignUpReferrerlReward create new object of signup referrer reward calculator type.
func NewSignUpReferrerReward() SignUpReferrerReward { return SignUpReferrerReward{} }

// ValidateCampaign validates campaign definitions.
func (sur SignUpReferrerReward) ValidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_REFERRAL {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "wrong reward category")
	}
	if campaign.RewardAmount.SubaccountAmount.LTE(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "signup rewards for subaccount should be positive")
	}
	if campaign.RewardAmountType != RewardAmountType_REWARD_AMOUNT_TYPE_FIXED {
		return sdkerrors.Wrapf(ErrWrongRewardAmountType, "reward amount type not supported for given reward type.")
	}

	return nil
}

// Calculate parses ticket payload and returns the distribution list of signup reward.
func (sur SignUpReferrerReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket, creator string,
) (RewardFactoryData, error) {
	var payload GrantSignupReferrerRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	if err := payload.Common.Validate(); err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "%s", err)
	}

	receiverAddr, err := sdk.AccAddressFromBech32(payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	promoter, isFound := keepers.GetPromoterByAddress(ctx, campaign.Promoter)
	if !isFound {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "promoter with the address: %s not found", campaign.Promoter)
	}

	if !keepers.RewardKeeper.HasRewardOfReceiverByPromoter(ctx, promoter.PromoterUID, payload.Referee, RewardCategory_REWARD_CATEGORY_SIGNUP) {
		return RewardFactoryData{}, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "referee account has signed up yet, there is no referee claim record")
	}

	if keepers.SubaccountKeeper.IsSubaccount(ctx, receiverAddr) {
		return RewardFactoryData{}, ErrReceiverAddrCanNotBeSubaccount
	}

	subaccountAddrStr, err := keepers.getSubaccountAddr(ctx, creator, payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	return NewRewardFactoryData(
		NewReceiver(
			payload.Common.Receiver,
			subaccountAddrStr,
			campaign.RewardAmount.MainAccountAmount,
			campaign.RewardAmount.SubaccountAmount,
			sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(),
			campaign.RewardAmount.UnlockPeriod,
		),
		payload.Common,
	), nil
}
