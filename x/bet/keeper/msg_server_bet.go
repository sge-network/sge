package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
)

func (k msgServer) PlaceBet(goCtx context.Context, msg *types.MsgPlaceBet) (*types.MsgPlaceBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetBetID(ctx, msg.Bet.UID)
	if isFound {
		return nil, sdkerrors.Wrapf(types.ErrDuplicateUID, "%s", msg.Bet.UID)
	}

	ticketData := &types.BetPlacementTicketPayload{}
	err := k.dvmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Bet.Ticket, &ticketData)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err = ticketData.Validate(msg.Creator); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketValidation, "%s", err)
	}

	bet := types.NewBet(msg.Creator, msg.Bet, ticketData.OddsType, ticketData.SelectedOdds)

	if err := k.Keeper.PlaceBet(ctx, bet); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInBetPlacement, "%s", err)
	}

	emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator)

	return &types.MsgPlaceBetResponse{Bet: msg.Bet}, nil
}

func emitBetEvent(ctx sdk.Context, msgType string, betUID string, betCreator string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msgType,
			sdk.NewAttribute(types.AttributeKeyBetCreator, betCreator),
			sdk.NewAttribute(types.AttributeKeyBetUID, betUID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msgType),
			sdk.NewAttribute(sdk.AttributeKeySender, betCreator),
		),
	})
}
