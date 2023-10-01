package keeper

import (
	"context"

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

	if msg.Creator != payload.FunderAddress {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, payload.FunderAddress, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	if err := payload.Validate(uint64(ctx.BlockTime().Unix())); err != nil {
		return nil, err
	}

	campaign := types.NewCampaign(
		msg.Creator, payload.FunderAddress, msg.Uid,
		payload.StartTs, payload.EndTs,
		payload.Type,
		payload.RewardDefs,
		types.NewPool(payload.PoolAmount),
	)

	rewardFactory, err := campaign.GetRewardsFactory()
	if err != nil {
		return nil, err
	}

	for _, d := range campaign.RewardDefs {
		if err := d.ValidateBasic(uint64(ctx.BlockTime().Unix())); err != nil {
			return nil, err
		}
	}

	err = rewardFactory.VaidateDefinitions(campaign)
	if err != nil {
		return nil, err
	}

	// transfer the pool amount to the reward pool module account
	if err := k.modFunder.Fund(
		types.RewardPoolFunder{}, ctx,
		sdk.MustAccAddressFromBech32(payload.FunderAddress),
		payload.PoolAmount,
	); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInFundingCampaignPool, "%s", err)
	}

	k.SetCampaign(ctx, campaign)

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

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.FunderAddress {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, valFound.FunderAddress, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	valFound.EndTS = payload.EndTs

	k.SetCampaign(ctx, valFound)

	return &types.MsgUpdateCampaignResponse{}, nil
}
