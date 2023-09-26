package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) ApplyReward(goCtx context.Context, msg *types.MsgApplyReward) (*types.MsgApplyRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	campaign, isFound := k.GetCampaign(ctx, msg.CampaignUid)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "campaign with the same uid is already set")
	}

	if err := campaign.CheckExpiration(uint64(ctx.BlockTime().Unix())); err == nil {
		return nil, err
	}

	rewardFactory, err := campaign.GetRewardsFactory()
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "failed to retrieve reward factory")
	}

	distribution, err := rewardFactory.CalculateDistributions(goCtx, ctx,
		types.RewardFactoryKeepers{
			OVMKeeper: k.ovmKeeper,
			BetKeeper: k.betKeeper,
		},
		campaign.RewardDefs, msg.Ticket)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "failed to get destribution from the ticket")
	}

	if err := campaign.CheckPoolBalance(distribution); err != nil {
		return nil, types.ErrInsufficientPoolBalance
	}

	if err := k.DistributeRewards(ctx, campaign.FunderAddress, distribution); err != nil {
		return nil, types.ErrInDistributionOfRewards
	}

	k.UpdateCampaignPool(ctx, campaign, distribution)

	return &types.MsgApplyRewardResponse{}, nil
}
