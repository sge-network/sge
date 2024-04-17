package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	bettypes "github.com/sge-network/sge/x/bet/types"
)

// BetBonusReward is the type for bet bonus rewards calculations
type BetBonusReward struct{}

// NewBetBonusReward create new object of bet bonus reward calculator type.
func NewBetBonusReward() BetBonusReward { return BetBonusReward{} }

// ValidateCampaign validates campaign definitions.
func (sur BetBonusReward) ValidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_BET_DISCOUNT {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "bet bonus rewards can only have single definition")
	}
	if campaign.RewardAmount.MainAccountPercentage.IsZero() && campaign.RewardAmount.SubaccountPercentage.IsZero() {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "one of main account and sub account percentage should be higher than zero")
	}
	if campaign.RewardAmountType != RewardAmountType_REWARD_AMOUNT_TYPE_PERCENTAGE {
		return sdkerrors.Wrapf(ErrWrongRewardAmountType, "reward amount type not supported for given reward type.")
	}
	if campaign.Constraints == nil || campaign.Constraints.MaxBetAmount.IsNil() {
		return sdkerrors.Wrapf(ErrMissingConstraintForCampaign, "constraints and max bet amount should be set for the bet bonus reward.")
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

	if keepers.SubaccountKeeper.IsSubaccount(ctx, addr) {
		return RewardFactoryData{}, ErrReceiverAddrCanNotBeSubaccount
	}

	subaccountAddrStr, err := keepers.getSubaccountAddr(ctx, creator, payload.Common.Receiver)
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

	if !bet.Meta.IsMainMarket {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "bet bonus grant is allowed for main market bets only %s", payload.BetUID)
	}

	if bet.Result != bettypes.Bet_RESULT_LOST &&
		bet.Result != bettypes.Bet_RESULT_WON {
		return RewardFactoryData{}, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "bet should be winner or loser, requested bet result is %s", bet.Result)
	}

	effectiveBetAmount := sdk.NewDecFromInt(bet.Amount)
	if campaign.Constraints != nil {
		if !campaign.Constraints.MaxBetAmount.IsNil() && campaign.Constraints.MaxBetAmount.GT(sdkmath.ZeroInt()) {
			effectiveBetAmount = sdk.NewDecFromInt(
				sdkmath.MinInt(campaign.Constraints.MaxBetAmount, bet.Amount),
			)
		}
	}

	mainAmount := effectiveBetAmount.Mul(campaign.RewardAmount.MainAccountPercentage).TruncateInt()
	subAmount := effectiveBetAmount.Mul(campaign.RewardAmount.SubaccountPercentage).TruncateInt()

	return NewRewardFactoryData(
		NewReceiver(
			payload.Common.Receiver,
			subaccountAddrStr,
			mainAmount, subAmount,
			campaign.RewardAmount.MainAccountPercentage, campaign.RewardAmount.SubaccountPercentage,
			campaign.RewardAmount.UnlockPeriod,
		),
		payload.Common,
	), nil
}
