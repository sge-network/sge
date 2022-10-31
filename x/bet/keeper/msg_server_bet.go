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
	_, isFound := k.GetBet(ctx, msg.Bet.UID)
	if isFound {
		emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator, types.AttributeValueStatusFailed, types.ErrDuplicateUID.Error())
		return &types.MsgPlaceBetResponse{}, nil
	}

	ticketData := &types.BetOdds{}
	err := k.dvmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Bet.Ticket, &ticketData)
	if err != nil {
		emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator, types.AttributeValueStatusFailed, sdkerrors.Wrapf(types.ErrInVerification, "%s", err).Error())
		return &types.MsgPlaceBetResponse{}, nil
	}

	if err = types.TicketFieldsValidation(ticketData); err != nil {
		emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator, types.AttributeValueStatusFailed, err.Error())
		return &types.MsgPlaceBetResponse{}, nil
	}

	bet, err := types.NewBet(msg.Creator, msg.Bet, ticketData)
	if err != nil {
		emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator, types.AttributeValueStatusFailed, err.Error())
		return &types.MsgPlaceBetResponse{}, nil
	}

	if err := k.Keeper.PlaceBet(ctx, bet); err != nil {
		emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator, types.AttributeValueStatusFailed, err.Error())
		return &types.MsgPlaceBetResponse{}, nil
	}

	emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator, types.AttributeValueStatusSuccessful, "")
	return &types.MsgPlaceBetResponse{}, nil
}

func (k msgServer) SettleBet(goCtx context.Context, msg *types.MsgSettleBet) (*types.MsgSettleBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.SettleBet(ctx, msg.BetUID); err != nil {
		emitBetEvent(ctx, types.TypeMsgSettleBet, msg.BetUID, msg.Creator, types.AttributeValueStatusFailed, err.Error())
		return &types.MsgSettleBetResponse{}, nil
	}
	emitBetEvent(ctx, types.TypeMsgSettleBet, msg.BetUID, msg.Creator, types.AttributeValueStatusSuccessful, "")
	return &types.MsgSettleBetResponse{}, nil
}

func emitBetEvent(ctx sdk.Context, msgType string, betUid string, betCreator string, status string, errorMsg string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msgType,
			sdk.NewAttribute(types.AttributeKeyBetCreator, betCreator),
			sdk.NewAttribute(types.AttributeKeyBetUID, betUid),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msgType),
			sdk.NewAttribute(sdk.AttributeKeySender, betCreator),
			sdk.NewAttribute(types.AttributeKeyStatus, status),
			sdk.NewAttribute(types.AttributeKeyErrorMessage, errorMsg),
		),
	})
}
