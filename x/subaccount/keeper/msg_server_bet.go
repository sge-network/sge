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

func (m msgServer) Wager(goCtx context.Context, msg *types.MsgWager) (*types.MsgWagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// find subaccount
	subAccountAddress, exists := m.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(msg.Msg.Creator))
	if !exists {
		return nil, status.Error(codes.NotFound, "subaccount not found")
	}

	// TODO: duplicate code from x/bet/keeper/msg_server_bet.go

	// Check if the value already exists
	_, isFound := m.keeper.betKeeper.GetBetID(ctx, msg.Msg.Props.UID)
	if isFound {
		return nil, sdkerrors.Wrapf(bettypes.ErrDuplicateUID, "%s", msg.Msg.Props.UID)
	}

	payload := &bettypes.WagerTicketPayload{}
	err := m.keeper.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Msg.Props.Ticket, &payload)
	if err != nil {
		return nil, sdkerrors.Wrapf(bettypes.ErrInTicketVerification, "%s", err)
	}

	originalSender := msg.Msg.Creator

	if err = payload.Validate(originalSender); err != nil {
		return nil, sdkerrors.Wrapf(bettypes.ErrInTicketValidation, "%s", err)
	}

	// duplication end

	// here we swap the original sender with the subaccount address
	bet := bettypes.NewBet(subAccountAddress.String(), msg.Msg.Props, payload.OddsType, payload.SelectedOdds)

	// make subaccount balance adjustments
	balance, exists := m.keeper.GetBalance(ctx, subAccountAddress)
	if !exists {
		panic("state corruption: subaccount balance not found")
	}

	err = balance.Spend(bet.Amount)
	if err != nil {
		return nil, err
	}

	if err := m.keeper.betKeeper.Wager(ctx, bet); err != nil {
		return nil, sdkerrors.Wrapf(bettypes.ErrInWager, "%s", err)
	}

	m.keeper.SetBalance(ctx, subAccountAddress, balance)

	msg.Msg.EmitEvent(&ctx)

	return &types.MsgWagerResponse{
		Response: &bettypes.MsgWagerResponse{Props: msg.Msg.Props},
	}, nil
}
