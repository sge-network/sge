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
func (sur SignUpRefereelReward) ValidateCampaign(campaign Campaign, blockTime uint64) error {
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

	sourceUIDAddr, err := sdk.AccAddressFromBech32(payload.Common.SourceUID)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "source address is invalid %s", err)
	}

	if acc := keepers.AccountKeeper.GetAccount(ctx, sourceUIDAddr); acc == nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "source id is not a valid account %s", payload.Common.SourceUID)
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
