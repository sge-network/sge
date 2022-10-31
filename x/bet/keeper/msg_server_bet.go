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
		return nil, types.ErrDuplicateUID
	}

	ticketData := &types.BetOdds{}
	err := k.dvmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Bet.Ticket, &ticketData)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	if err = types.TicketFieldsValidation(ticketData); err != nil {
		return nil, err
	}

	bet, err := types.NewBet(msg.Creator, msg.Bet, ticketData)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.PlaceBet(ctx, bet); err != nil {
		return nil, err
	}
	return &types.MsgPlaceBetResponse{}, nil
}

func (k msgServer) SettleBet(goCtx context.Context, msg *types.MsgSettleBet) (*types.MsgSettleBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.SettleBet(ctx, msg.BetUID); err != nil {
		return nil, err
	}
	return &types.MsgSettleBetResponse{}, nil
}
