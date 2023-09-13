package keeper

import (
	"context"

	cosmerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	campaign := types.Campaign{
		Creator: msg.Creator,
		UID:     msg.Uid,
	}

	k.SetCampaign(
		ctx,
		campaign,
	)
	return &types.MsgCreateCampaignResponse{}, nil
}

func (k msgServer) UpdateCampaign(goCtx context.Context, msg *types.MsgUpdateCampaign) (*types.MsgUpdateCampaignResponse, error) {
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
	if msg.Creator != valFound.Creator {
		return nil, cosmerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	campaign := types.Campaign{
		Creator: msg.Creator,
		UID:     msg.Uid,
	}

	k.SetCampaign(ctx, campaign)

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
	if msg.Creator != valFound.Creator {
		return nil, cosmerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveCampaign(
		ctx,
		msg.Uid,
	)

	return &types.MsgDeleteCampaignResponse{}, nil
}
