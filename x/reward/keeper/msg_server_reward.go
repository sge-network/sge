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

	recevier, rewardCommon, isSubAccount, oneTimeKey, err := rewardFactory.Calculate(goCtx, ctx,
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
	k.SetReward(ctx, types.NewReward(
		msg.Uid, msg.Creator, recevier.Addr,
		msg.CampaignUid, campaign.RewardAmount,
		rewardCommon.Source, rewardCommon.SourceCode, rewardCommon.SourceUID,
		uint64(ctx.BlockTime().Unix()),
	))
	k.SetOneTimeReward(ctx, types.NewOneTimeReward(oneTimeKey, campaign.RewardType))
	k.SetRewardByReceiver(ctx, types.NewRewardByType(msg.Uid, recevier.Addr, campaign.RewardType))
	k.SetRewardByCampaign(ctx, types.NewRewardByCampaign(msg.Uid, campaign.UID))

	msg.EmitEvent(&ctx, msg.CampaignUid, recevier)

	return &types.MsgGrantRewardResponse{}, nil
}
