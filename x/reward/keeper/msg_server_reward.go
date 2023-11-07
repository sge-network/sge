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

	rewards, err := k.GetRewardsByAddressAndCategory(ctx, msg.Creator, campaign.RewardCategory)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "failed to retrieve rewards for user.")
	}
	if len(rewards) > 0 {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "maximum rewards claimed for the given category.")
	}

	recevier, rewardCommon, isSubAccount, _, err := rewardFactory.Calculate(goCtx, ctx,
		types.RewardFactoryKeepers{
			OVMKeeper:        k.ovmKeeper,
			BetKeeper:        k.betKeeper,
			SubAccountKeeper: k.subaccountKeeper,
		}, campaign, msg.Ticket)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "distribution calculation failed %s", err)
	}

	if err := campaign.CheckPoolBalance(recevier.Amount); err != nil {
		return nil, types.ErrInsufficientPoolBalance
	}

	if err := k.DistributeRewards(ctx, campaign.Promoter, isSubAccount, recevier); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInDistributionOfRewards, "%s", err)
	}

	k.UpdateCampaignPool(ctx, campaign, recevier)
	reward := types.NewReward(
		msg.Uid, msg.Creator, recevier.Addr,
		msg.CampaignUid, campaign.RewardAmount,
		rewardCommon.SourceUID,
		"",
	)
	k.SetReward(ctx, reward)
	k.SetRewardByReceiver(ctx, campaign.RewardCategory, reward)
	k.SetRewardByCampaign(ctx, reward)

	msg.EmitEvent(&ctx, msg.CampaignUid, recevier)

	return &types.MsgGrantRewardResponse{}, nil
}
