package keeper

import (
	"context"

	cosmerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) CreateCampaign(goCtx context.Context, msg *types.MsgCreateCampaign) (*types.MsgCreateCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetCampaign(
		ctx,
		msg.Uid,
	)
	if isFound {
		return nil, cosmerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var payload types.CreateCampaignPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, cosmerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if msg.Creator != payload.FunderAddress {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, payload.FunderAddress, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	// transfer the pool amount to the reward pool module account
	if err := k.modFunder.Fund(
		types.RewardPoolFunder{}, ctx,
		sdk.MustAccAddressFromBech32(payload.FunderAddress),
		payload.PoolAmount,
	); err != nil {
		return nil, cosmerrors.Wrapf(types.ErrInFundingCampaignPool, "%s", err)
	}

	campaign := types.NewCampaign(
		msg.Creator, payload.FunderAddress, msg.Uid,
		payload.StartTs, payload.EndTs,
		payload.Type,
		payload.RewardDefs,
		types.NewPool(payload.PoolAmount),
	)
	k.SetCampaign(ctx, campaign)

	return &types.MsgCreateCampaignResponse{}, nil
}

func (k msgServer) UpdateCampaign(goCtx context.Context, msg *types.MsgUpdateCampaign) (*types.MsgUpdateCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.UpdateCampaignPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, cosmerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	// Check if the value exists
	valFound, isFound := k.GetCampaign(ctx, msg.Uid)
	if !isFound {
		return nil, cosmerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
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

func (k msgServer) DeleteCampaign(goCtx context.Context, msg *types.MsgDeleteCampaign) (*types.MsgDeleteCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetCampaign(
		ctx,
		msg.Uid,
	)
	if !isFound {
		return nil, cosmerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.FunderAddress {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, valFound.FunderAddress, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	k.RemoveCampaign(ctx, msg.Uid)

	return &types.MsgDeleteCampaignResponse{}, nil
}
