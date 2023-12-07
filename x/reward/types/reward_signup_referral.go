package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// SignUpReferralReward is the type for signup referral rewards calculations
type SignUpReferralReward struct{}

// NewSignUpReferralReward create new object of signup referral reward calculator type.
func NewSignUpReferralReward() SignUpReward { return SignUpReward{} }

// VaidateCampaign validates campaign definitions.
func (sur SignUpReferralReward) ValidateCampaign(campaign Campaign, blockTime uint64) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_SIGNUP {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "signup rewards can only have single definition")
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
func (sur SignUpReferralReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket, creator string,
) (RewardFactoryData, error) {
	var payload GrantSignupReferreeRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	addr, err := sdk.AccAddressFromBech32(payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	pairs := []string{
		payload.Common.SourceUID, // referrer address
		payload.Common.Receiver,
	}

	if keepers.SubAccountKeeper.IsSubAccount(ctx, addr) {
		return RewardFactoryData{}, ErrReceiverAddrCanNotBeSubAcc
	}

	subAccountAddressString, err := keepers.getSubAccAddr(ctx, creator, payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	return NewRewardFactoryData(
		NewReceiver(
			subAccountAddressString,
			payload.Common.Receiver,
			campaign.RewardAmount.SubaccountAmount,
			campaign.RewardAmount.MainAccountAmount,
			campaign.RewardAmount.UnlockPeriod,
		),
		payload.Common,
		pairs...,
	), nil
}
