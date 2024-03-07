package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var percent = sdk.NewInt(100)

// BetBonusReward is the type for bet bonus rewards calculations
type BetBonusReward struct{}

// NewBetBonusReward create new object of bet bonus reward calculator type.
func NewBetBonusReward() BetBonusReward { return BetBonusReward{} }

// ValidateCampaign validates campaign definitions.
func (sur BetBonusReward) ValidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_BET_DISCOUNT {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "bet bonus rewards can only have single definition")
	}
	if campaign.RewardAmount.MainAccountAmount.GT(percent) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "bet bonus rewards percent for main account should be between 0 and 100")
	}
	if campaign.RewardAmount.SubaccountAmount.GT(percent) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "bet bonus rewards percent for sub account should be between 0 and 100")
	}
	if campaign.RewardAmount.MainAccountAmount.IsZero() && campaign.RewardAmount.SubaccountAmount.IsZero() {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "one of main account and sub account percentage should be higher than zero")
	}
	if campaign.RewardAmountType != RewardAmountType_REWARD_AMOUNT_TYPE_PERCENTAGE {
		return sdkerrors.Wrapf(ErrWrongRewardAmountType, "reward amount type not supported for given reward type.")
	}

	return nil
}

// Calculate parses ticket payload and returns the distribution list of bet bonus reward.
func (sur BetBonusReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket, creator string,
) (RewardFactoryData, error) {
	var payload GrantBetBonusRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	if err := payload.Common.Validate(); err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "%s", err)
	}

	addr, err := sdk.AccAddressFromBech32(payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	if keepers.SubAccountKeeper.IsSubAccount(ctx, addr) {
		return RewardFactoryData{}, ErrReceiverAddrCanNotBeSubAcc
	}

	subAccountAddressString, err := keepers.getSubAccAddr(ctx, creator, payload.Common.Receiver)
	if err != nil {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	uid2ID, found := keepers.BetKeeper.GetBetID(ctx, payload.BetUID)
	if !found {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "bet id not found for uid %s", payload.BetUID)
	}

	bet, found := keepers.BetKeeper.GetBet(ctx, payload.Common.Receiver, uid2ID.ID)
	if !found {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "bet not found with uid %s", payload.BetUID)
	}

	mainAmount := bet.Amount.Mul(campaign.RewardAmount.MainAccountAmount).Quo(percent)
	subAmount := bet.Amount.Mul(campaign.RewardAmount.SubaccountAmount).Quo(percent)

	return NewRewardFactoryData(
		NewReceiver(
			subAccountAddressString,
			payload.Common.Receiver,
			mainAmount,
			subAmount,
			campaign.RewardAmount.UnlockPeriod,
		),
		payload.Common,
	), nil
}
