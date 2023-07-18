package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
)

func (k msgServer) Place(
	goCtx context.Context,
	msg *types.MsgPlace,
) (*types.MsgPlaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetBetID(ctx, msg.Bet.UID)
	if isFound {
		return nil, sdkerrors.Wrapf(types.ErrDuplicateUID, "%s", msg.Bet.UID)
	}

	payload := &types.PlacementTicketPayload{}
	err := k.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Bet.Ticket, &payload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err = payload.Validate(msg.Creator); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketValidation, "%s", err)
	}

	bet := types.NewBet(msg.Creator, msg.Bet, payload.OddsType, payload.SelectedOdds)

	if err := k.Keeper.Place(ctx, bet); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInBetPlacement, "%s", err)
	}

	msg.EmitEvent(&ctx)

	return &types.MsgPlaceResponse{Bet: msg.Bet}, nil
}
