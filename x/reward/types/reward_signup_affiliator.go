package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// SignUpAffiliatorReward is the type for signup Affiliator rewards calculations
type SignUpAffiliatorReward struct{}

// NewSignUpAffiliatorlReward create new object of signup Affiliator reward calculator type.
func NewSignUpAffiliatorReward() SignUpAffiliatorReward { return SignUpAffiliatorReward{} }

// ValidateCampaign validates campaign definitions.
func (sur SignUpAffiliatorReward) ValidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_AFFILIATE {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "wrong reward category")
	}
	if campaign.RewardAmount.SubaccountAmount.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "affiliate rewards for subaccount should be zero")
	}
	if campaign.RewardAmount.MainAccountAmount.LTE(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "affiliate rewards for main account should be positive")
	}
	if campaign.RewardAmountType != RewardAmountType_REWARD_AMOUNT_TYPE_FIXED {
		return sdkerrors.Wrapf(ErrWrongRewardAmountType, "reward amount type not supported for given reward type.")
	}

	return nil
}

// Calculate parses ticket payload and returns the distribution list of signup reward.
func (sur SignUpAffiliatorReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket, creator string,
) (RewardFactoryData, error) {
	var payload GrantSignupAffiliatorRewardPayload
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
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "promoter with the address: %s not found", &campaign.Promoter)
	}

	if !keepers.RewardKeeper.HasRewardOfReceiverByPromoter(ctx, promoter.PromoterUID, payload.Affiliatee, RewardCategory_REWARD_CATEGORY_SIGNUP) {
		return RewardFactoryData{}, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "affiliatee account has signed up yet, there is no affiliatee claim record")
	}

	if keepers.SubAccountKeeper.IsSubAccount(ctx, receiverAddr) {
		return RewardFactoryData{}, ErrReceiverAddrCanNotBeSubAcc
	}

	subAccAddrStr, err := keepers.getSubAccAddr(ctx, creator, payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	return NewRewardFactoryData(
		NewReceiver(
			subAccAddrStr,
			payload.Common.Receiver,
			campaign.RewardAmount.SubaccountAmount,
			campaign.RewardAmount.MainAccountAmount,
			campaign.RewardAmount.UnlockPeriod,
		),
		payload.Common,
	), nil
}
