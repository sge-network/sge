package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) ApplyReward(goCtx context.Context, msg *types.MsgApplyReward) (*types.MsgApplyRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	campaign, isFound := k.GetCampaign(ctx, msg.CampaignUid)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrInvalidRequest, "campaign with the same uid is already set")
	}

	rewardFactory, err := campaign.GetRewardsFactory()
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrInvalidRequest, "failed to retrieve reward factory")
	}

	if err := rewardFactory.ValidateBasic(campaign); err == nil {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrInvalidRequest, "basic validation failed")
	}

	distribution, err := rewardFactory.CalculateDistributions(campaign.RewardDefs, msg.Ticket)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrInvalidRequest, "failed to get destribution from the ticket")
	}

	if err := campaign.CheckPoolBalance(distribution); err != nil {
		return nil, types.ErrInsufficientPoolBalance
	}

	if err := k.DistributeRewards(ctx, distribution); err != nil {
		return nil, types.ErrInDistributionOfRewards
	}

	k.UpdateCampaignPool(ctx, campaign, distribution)

	return &types.MsgApplyRewardResponse{}, nil
}
