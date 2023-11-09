package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) CreateCampaign(goCtx context.Context, msg *types.MsgCreateCampaign) (*types.MsgCreateCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetCampaign(ctx, msg.Uid)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, "index already set")
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

	if err := payload.Validate(uint64(ctx.BlockTime().Unix())); err != nil {
		return nil, err
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
		types.NewPool(payload.TotalFunds),
	)

	rewardFactory, err := campaign.GetRewardsFactory()
	if err != nil {
		return nil, err
	}

	err = rewardFactory.VaidateCampaign(campaign, uint64(ctx.BlockTime().Unix()))
	if err != nil {
		return nil, err
	}

	// transfer the pool amount to the reward pool module account
	if err := k.modFunder.Fund(
		types.RewardPoolFunder{}, ctx,
		sdk.MustAccAddressFromBech32(payload.Promoter),
		payload.TotalFunds,
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

	if err := payload.Validate(uint64(ctx.BlockTime().Unix())); err != nil {
		return nil, err
	}

	// Check if the value exists
	valFound, isFound := k.GetCampaign(ctx, msg.Uid)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrtypes.ErrKeyNotFound, "index not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != valFound.Promoter {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, valFound.Promoter, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	valFound.EndTS = payload.EndTs

	k.SetCampaign(ctx, valFound)

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
