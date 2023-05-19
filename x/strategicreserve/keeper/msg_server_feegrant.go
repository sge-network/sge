package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// InvokeFeeGrant accepts ticket containing creation feegrant and return response after processing
func (k msgServer) InvokeFeeGrant(goCtx context.Context, msg *types.MsgInvokeFeeGrant) (*types.MsgInvokeFeeGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.InvokeFeeGrantPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err := payload.Validate(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	_, found := k.Keeper.GetFeeGrant(ctx, payload.Grantee)
	if found {
		return nil, types.ErrFeeGrantExists
	}

	grantee, err := sdk.AccAddressFromBech32(payload.Grantee)
	if err != nil {
		return nil, err
	}

	err = k.SetSrPoolSdkFeeGrant(ctx, msg.Creator, grantee)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInSrPoolFeeGrant, "%s", err)
	}

	feeGrant := types.NewFeeGrant(
		msg.Creator,
		grantee.String(),
		ctx.BlockTime().Unix(),
	)

	k.Keeper.SetFeeGrant(ctx, feeGrant)

	emitFeeGrantEvent(ctx, types.TypeMsgInvokeFeeGrant, msg.Creator, feeGrant.Grantee)

	return &types.MsgInvokeFeeGrantResponse{
		FeeGrant: &feeGrant,
	}, nil
}

// RevokeFeeGrant accepts ticket containing revoke feegrant and return response after processing.
func (k msgServer) RevokeFeeGrant(goCtx context.Context, msg *types.MsgRevokeFeeGrant) (*types.MsgRevokeFeeGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.RevokeFeeGrantPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	feeGrant, found := k.Keeper.GetFeeGrant(ctx, payload.Grantee)
	if !found {
		return nil, types.ErrMarketNotFound
	}

	// update market is not valid, return error
	if err := payload.Validate(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	// update market is successful, update the module state
	k.Keeper.RemoveFeeGrant(ctx, feeGrant)

	emitFeeGrantEvent(ctx, types.TypeMsgRevokeFeeGrant, msg.Creator, feeGrant.Grantee)

	return &types.MsgRevokeFeeGrantResponse{FeeGrant: nil}, nil
}

func emitFeeGrantEvent(ctx sdk.Context, emitType string, creator string, grantee string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			emitType,
			sdk.NewAttribute(types.AttributeKeyFeeGrantCreator, creator),
			sdk.NewAttribute(types.AttributeKeyGrantee, grantee),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, emitType),
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
		),
	})
}
