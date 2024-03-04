package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) CreateCampaign(goCtx context.Context, msg *types.MsgCreateCampaign) (*types.MsgCreateCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetCampaign(ctx, msg.Uid)
	if isFound {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "campaign with the uid: %s already exists", msg.Uid)
	}

	var payload types.CreateCampaignPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if msg.Creator != payload.Promoter {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, payload.Promoter, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	if err := payload.Validate(cast.ToUint64(ctx.BlockTime().Unix())); err != nil {
		return nil, err
	}

	totalRewardAmount := sdkmath.ZeroInt()
	if !payload.RewardAmount.MainAccountAmount.IsNil() {
		totalRewardAmount = totalRewardAmount.Add(payload.RewardAmount.MainAccountAmount)
	}
	if !payload.RewardAmount.SubaccountAmount.IsNil() {
		totalRewardAmount = totalRewardAmount.Add(payload.RewardAmount.SubaccountAmount)
	}

	if msg.TotalFunds.LT(totalRewardAmount) {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "defined reward amount %s is more than total funds %s", totalRewardAmount, msg.TotalFunds)
	}

	campaign := types.NewCampaign(
		msg.Creator, payload.Promoter, msg.Uid,
		payload.StartTs, payload.EndTs, payload.ClaimsPerCategory,
		payload.RewardType,
		payload.Category,
		payload.RewardAmountType,
		payload.RewardAmount,
		payload.IsActive,
		payload.Meta,
		types.NewPool(msg.TotalFunds),
	)

	rewardFactory, err := campaign.GetRewardsFactory()
	if err != nil {
		return nil, err
	}

	err = rewardFactory.ValidateCampaign(campaign)
	if err != nil {
		return nil, err
	}

	// transfer the pool amount to the reward pool module account
	if err := k.modFunder.Fund(
		types.RewardPoolFunder{}, ctx,
		sdk.MustAccAddressFromBech32(payload.Promoter),
		msg.TotalFunds,
	); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInFundingCampaignPool, "%s", err)
	}

	k.SetCampaign(ctx, campaign)

	msg.EmitEvent(&ctx, msg.Uid)

	return &types.MsgCreateCampaignResponse{}, nil
}

func (k msgServer) UpdateCampaign(goCtx context.Context, msg *types.MsgUpdateCampaign) (*types.MsgUpdateCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.UpdateCampaignPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err := payload.Validate(cast.ToUint64(ctx.BlockTime().Unix())); err != nil {
		return nil, err
	}

	// Check if the value exists
	campaign, isFound := k.GetCampaign(ctx, msg.Uid)
	if !isFound {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrKeyNotFound, "campaign with the id: %s does not exist", msg.Uid)
	}

	if !campaign.IsActive {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "inactive campaign can not be updated")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != campaign.Promoter {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, campaign.Promoter, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	if !msg.TopupFunds.IsNil() && msg.TopupFunds.GT(sdkmath.ZeroInt()) {
		// transfer the pool amount to the reward pool module account
		if err := k.modFunder.Fund(
			types.RewardPoolFunder{}, ctx,
			sdk.MustAccAddressFromBech32(campaign.Promoter),
			msg.TopupFunds,
		); err != nil {
			return nil, sdkerrors.Wrapf(types.ErrInFundingCampaignPool, "%s", err)
		}

		campaign.Pool.Total = campaign.Pool.Total.Add(msg.TopupFunds)
	}

	campaign.EndTS = payload.EndTs
	campaign.IsActive = payload.IsActive

	k.SetCampaign(ctx, campaign)

	msg.EmitEvent(&ctx, msg.Uid)

	return &types.MsgUpdateCampaignResponse{}, nil
}

func (k msgServer) WithdrawFunds(goCtx context.Context, msg *types.MsgWithdrawFunds) (*types.MsgWithdrawFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.WithdrawFundsPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	// Validate ticket payload
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// Check if the campaign exists
	valFound, isFound := k.GetCampaign(ctx, msg.Uid)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrKeyNotFound, "campaign not found")
	}

	if payload.Promoter != valFound.Promoter {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrKeyNotFound, "promoter should be the same as stored campaign promoter")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != valFound.Promoter {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, valFound.Promoter, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}
	availableAmount := valFound.Pool.Total.Sub(valFound.Pool.Spent)
	// check if the pool amount is positive
	if availableAmount.IsNil() || !availableAmount.GT(sdkmath.ZeroInt()) {
		return nil, sdkerrors.Wrapf(types.ErrWithdrawFromCampaignPool, "pool amount should be positive")
	}

	// transfer the funds present in campaign to the promoter
	if err := k.modFunder.Refund(
		types.RewardPoolFunder{}, ctx,
		sdk.MustAccAddressFromBech32(payload.Promoter),
		availableAmount,
	); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrWithdrawFromCampaignPool, "%s", err)
	}
	// set the pool amount to zero
	valFound.Pool.Total = sdkmath.ZeroInt()
	// deactivate the campaign
	valFound.IsActive = false

	// store the campaign
	k.SetCampaign(ctx, valFound)
	// emit withdraw event
	msg.EmitEvent(&ctx, msg.Uid)

	return &types.MsgWithdrawFundsResponse{}, nil
}
