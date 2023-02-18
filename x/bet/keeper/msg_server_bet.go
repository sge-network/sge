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
		return &types.MsgPlaceBetResponse{
				Error: types.ErrDuplicateUID.Error(),
				Bet:   msg.Bet,
			},
			types.ErrDuplicateUID
	}

	ticketData := &types.BetPlacementTicketPayload{}
	err := k.dvmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Bet.Ticket, &ticketData)
	if err != nil {
		return &types.MsgPlaceBetResponse{
				Error: sdkerrors.Wrapf(types.ErrInVerification, "%s", err).Error(),
				Bet:   msg.Bet,
			},
			sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	if err = types.TicketFieldsValidation(ticketData); err != nil {
		return &types.MsgPlaceBetResponse{
				Error: err.Error(),
				Bet:   msg.Bet,
			},
			err
	}

	// Kyc validation is done only when the KYC is required
	if !kycValidation(msg.Creator, ticketData.KycData) {
		return nil, sdkerrors.Wrapf(types.ErrUserKycFailed, "%s", msg.Creator)
	}

	bet, err := types.NewBet(msg.Creator, msg.Bet, ticketData.SelectedOdds)
	if err != nil {
		return &types.MsgPlaceBetResponse{
				Error: err.Error(),
				Bet:   msg.Bet,
			},
			err
	}

	if err := k.Keeper.PlaceBet(ctx, bet); err != nil {
		return &types.MsgPlaceBetResponse{
				Error: err.Error(),
				Bet:   msg.Bet,
			},
			err
	}

	emitBetEvent(ctx, types.TypeMsgPlaceBet, msg.Bet.UID, msg.Creator)
	return &types.MsgPlaceBetResponse{
			Error: "",
			Bet:   msg.Bet,
		},
		nil
}

func (k msgServer) SettleBet(goCtx context.Context, msg *types.MsgSettleBet) (*types.MsgSettleBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.SettleBet(ctx, msg.BettorAddress, msg.BetUID); err != nil {
		return &types.MsgSettleBetResponse{
			Error:  err.Error(),
			BetUID: msg.BetUID,
		}, err
	}
	emitBetEvent(ctx, types.TypeMsgSettleBet, msg.BetUID, msg.Creator)
	return &types.MsgSettleBetResponse{}, nil
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
