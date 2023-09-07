package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) ApplyReward(goCtx context.Context, msg *types.MsgApplyReward) (*types.MsgApplyRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	campaign, isFound := k.GetCampaign(
		ctx,
		msg.CampaignUid,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	rewardFactory, err := GetRewardsFactory(campaign.RewardType)
	if err != nil {
		return nil, err
	}

	if err := rewardFactory.ValidateBasic(campaign); err == nil {
		return nil, err
	}

	distribution, err := rewardFactory.CalculateDistributions(campaign.RewardDefs, msg.Ticket)
	if err != nil {
		return nil, err
	}

	// TODO: call sge distribution module hook to transfer funds

	k.UpdateCampaignPool(ctx, campaign, distribution)

	return &types.MsgApplyRewardResponse{}, nil
}
