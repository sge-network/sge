package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
)

// BetDiscountReward is the type for no loss bets rewards calculations
type BetDiscountReward struct{}

// NewBetDiscountReward create new object of no loss bets reward calculator type.
func NewBetDiscountReward() BetDiscountReward { return BetDiscountReward{} }

// VaidateCampaign validates campaign definitions.
func (afr BetDiscountReward) VaidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_BET_DISCOUNT {
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

// Calculate parses ticket payload and returns the distribution list of no loss bets reward.
func (afr BetDiscountReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket string,
) (Allocation, error) {
	var payload GrantBetDiscountRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return Allocation{}, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	_, found := keepers.SubAccountKeeper.GetSubAccountOwner(ctx, sdk.MustAccAddressFromBech32(payload.Common.Receiver))
	if !found {
		return Allocation{}, sdkerrors.Wrapf(ErrReceiverAddrNotSubAcc, "%s", &payload.Common.Receiver)
	}

	bettorAddr := payload.Common.Receiver
	for _, betUID := range payload.BetUids {
		uID2ID, found := keepers.BetKeeper.GetBetID(ctx, betUID)
		if !found {
			return Allocation{}, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "bet id not found %s", betUID)
		}
		bet, found := keepers.BetKeeper.GetBet(ctx, bettorAddr, uID2ID.ID)
		if !found {
			return Allocation{}, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "bet not found %s", betUID)
		}
		if bet.Result != bettypes.Bet_RESULT_LOST {
			return Allocation{}, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "the bet result is not loss %s", betUID)
		}
	}

	return Allocation{
		SubAcc: NewReceiver(
			payload.Common.Receiver,
			campaign.RewardAmount.SubaccountAmount,
			0,
		),
	}, nil
}
