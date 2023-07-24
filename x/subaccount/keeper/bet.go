package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/subaccount/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m msgServer) PlaceBet(goCtx context.Context, msg *types.MsgPlaceBet) (*types.MsgPlaceBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// find subaccount
	subAccountAddress, exists := m.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(msg.Msg.Creator))
	if !exists {
		return nil, status.Error(codes.NotFound, "subaccount not found")
	}

	// TODO: duplicate code from x/bet/keeper/msg_server_bet.go

	// Check if the value already exists
	_, isFound := m.keeper.betKeeper.GetBetID(ctx, msg.Msg.Bet.UID)
	if isFound {
		return nil, sdkerrors.Wrapf(bettypes.ErrDuplicateUID, "%s", msg.Msg.Bet.UID)
	}

	payload := &bettypes.BetPlacementTicketPayload{}
	err := m.keeper.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Msg.Bet.Ticket, &payload)
	if err != nil {
		return nil, sdkerrors.Wrapf(bettypes.ErrInTicketVerification, "%s", err)
	}

	originalSender := msg.Msg.Creator

	if err = payload.Validate(originalSender); err != nil {
		return nil, sdkerrors.Wrapf(bettypes.ErrInTicketValidation, "%s", err)
	}

	// here we swap the original sender with the subaccount address
	bet := bettypes.NewBet(subAccountAddress.String(), msg.Msg.Bet, payload.OddsType, payload.SelectedOdds)

	if err := m.keeper.betKeeper.PlaceBet(ctx, bet); err != nil {
		return nil, sdkerrors.Wrapf(bettypes.ErrInBetPlacement, "%s", err)
	}

	msg.Msg.EmitEvent(&ctx)

	return &types.MsgPlaceBetResponse{
		Response: &bettypes.MsgPlaceBetResponse{Bet: msg.Msg.Bet},
	}, nil

}
