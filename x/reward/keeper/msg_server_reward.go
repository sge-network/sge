package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) GrantReward(goCtx context.Context, msg *types.MsgGrantReward) (*types.MsgGrantRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, isFound := k.GetReward(ctx, msg.Uid); isFound {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "reward grant with uid: %s exists", msg.Uid)
	}

	campaign, isFound := k.GetCampaign(ctx, msg.CampaignUid)
	if !isFound {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "campaign with the uid: %s not found", msg.CampaignUid)
	}

	if !campaign.IsActive {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "campaign with uid: %s not active", msg.CampaignUid)
	}

	if err := campaign.CheckTS(cast.ToUint64(ctx.BlockTime().Unix())); err != nil {
		return nil, err
	}

	rewardFactory, err := campaign.GetRewardsFactory()
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "failed to retrieve reward factory")
	}

	factData, err := rewardFactory.Calculate(goCtx, ctx,
		types.RewardFactoryKeepers{
			OVMKeeper:        k.ovmKeeper,
			BetKeeper:        k.betKeeper,
			SubaccountKeeper: k.subaccountKeeper,
			RewardKeeper:     k.Keeper,
			AccountKeeper:    k.accountKeeper,
		}, campaign, msg.Ticket, msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "distribution calculation failed %s", err)
	}

	rewards, err := k.GetRewardsByAddressAndCategory(ctx, factData.Receiver.MainAccountAddr, campaign.RewardCategory)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "failed to retrieve rewards for user.")
	}
	if len(rewards) >= cast.ToInt(campaign.ClaimsPerCategory) {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "maximum rewards claimed for the given category.")
	}

	if err := campaign.CheckPoolBalance(factData.Receiver.SubaccountAmount.Add(factData.Receiver.MainAccountAmount)); err != nil {
		return nil, types.ErrInsufficientPoolBalance
	}

	unlockTS, err := k.DistributeRewards(ctx, factData.Receiver)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInDistributionOfRewards, "%s", err)
	}

	k.UpdateCampaignPool(ctx, campaign, factData.Receiver)
	k.SetReward(ctx, types.NewReward(
		msg.Uid, msg.Creator, factData.Receiver.MainAccountAddr,
		msg.CampaignUid, campaign.RewardAmount,
		factData.Common.SourceUID,
		factData.Common.Meta,
	))
	k.SetRewardByReceiver(ctx, types.NewRewardByCategory(msg.Uid, factData.Receiver.MainAccountAddr, campaign.RewardCategory))
	k.SetRewardByCampaign(ctx, types.NewRewardByCampaign(msg.Uid, campaign.UID))

	msg.EmitEvent(&ctx, msg.CampaignUid, msg.Uid, campaign.Promoter, *campaign.RewardAmount, factData.Receiver, unlockTS)

	return &types.MsgGrantRewardResponse{}, nil
}
