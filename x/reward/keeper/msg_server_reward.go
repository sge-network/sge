package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) GrantReward(goCtx context.Context, msg *types.MsgGrantReward) (*types.MsgGrantRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	campaign, isFound := k.GetCampaign(ctx, msg.CampaignUid)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "campaign with the uid not found")
	}

	if err := campaign.CheckExpiration(uint64(ctx.BlockTime().Unix())); err != nil {
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
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "distribution calculation failed %s", err)
	}

	if err := campaign.CheckPoolBalance(distribution); err != nil {
		return nil, types.ErrInsufficientPoolBalance
	}

	if err := k.DistributeRewards(ctx, campaign.Promoter, distribution); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInDistributionOfRewards, "%s", err)
	}

	k.UpdateCampaignPool(ctx, campaign, distribution)

	msg.EmitEvent(&ctx, msg.CampaignUid, distribution)

	return &types.MsgGrantRewardResponse{}, nil
}
