package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// SignUpRefereelReward is the type for signup referral rewards calculations
type SignUpRefereelReward struct{}

// NewSignUpRefereelReward create new object of signup referral reward calculator type.
func NewSignUpRefereelReward() SignUpRefereelReward { return SignUpRefereelReward{} }

// ValidateCampaign validates campaign definitions.
func (sur SignUpRefereelReward) ValidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_SIGNUP {
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
func (sur SignUpRefereelReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket, creator string,
) (RewardFactoryData, error) {
	var payload GrantSignupRewardPayload
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

	if _, err = sdk.AccAddressFromBech32(payload.Common.SourceUID); err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "source address is invalid %s", err)
	}

	if keepers.SubaccountKeeper.IsSubaccount(ctx, receiverAddr) {
		return RewardFactoryData{}, ErrReceiverAddrCanNotBeSubaccount
	}

	subaccountAddrStr, err := keepers.getSubaccountAddr(ctx, creator, payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	hasBet, err := keepers.BetKeeper.IsAnyBetForAccount(ctx, payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrPanic, "%s", err)
	}
	if !hasBet {
		return RewardFactoryData{}, ErrNoBetForReceiverFound
	}

	return NewRewardFactoryData(
		NewReceiver(
			payload.Common.Receiver,
			subaccountAddrStr,
			campaign.RewardAmount.MainAccountAmount,
			campaign.RewardAmount.SubaccountAmount,
			sdk.ZeroDec(), sdk.ZeroDec(),
			campaign.RewardAmount.UnlockPeriod,
		),
		payload.Common,
	), nil
}
