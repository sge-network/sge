package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/reward/types"
)

func (k msgServer) CreatePromoter(goCtx context.Context, msg *types.MsgCreatePromoter) (*types.MsgCreatePromoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.CreatePromoterPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	// Check if the value already exists
	_, isFound := k.GetPromoter(ctx, payload.UID)
	if isFound {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "promoter with the uid: %s already exists", payload.UID)
	}

	if err := payload.Validate(); err != nil {
		return nil, err
	}

	k.SetPromoter(ctx, types.Promoter{
		Creator:   msg.Creator,
		Addresses: []string{msg.Creator},
		UID:       payload.UID,
		Conf:      payload.Conf,
	})

	msg.EmitEvent(&ctx, payload.UID, payload.Conf)

	return &types.MsgCreatePromoterResponse{}, nil
}

func (k msgServer) SetPromoterConf(goCtx context.Context, msg *types.MsgSetPromoterConf) (*types.MsgSetPromoterConfResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	promoter, isFound := k.GetPromoter(ctx, msg.Uid)
	if !isFound {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "promoter does not exist %s", msg.Uid)
	}

	creatorIsPromoter := false
	for _, p := range promoter.Addresses {
		if p == msg.Creator {
			creatorIsPromoter = true
		}
	}

	if !creatorIsPromoter {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "creator should be one of stored addresses in promoter")
	}

	var payload types.SetPromoterConfPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err := payload.Validate(); err != nil {
		return nil, err
	}

	promoter.Conf = payload.Conf

	k.SetPromoter(ctx, promoter)

	msg.EmitEvent(&ctx, promoter.Conf)

	return &types.MsgSetPromoterConfResponse{}, nil
}
