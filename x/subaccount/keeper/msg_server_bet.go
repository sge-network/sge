package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (k msgServer) Wager(goCtx context.Context, msg *types.MsgWager) (*types.MsgWagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	subAccOwner := sdk.MustAccAddressFromBech32(msg.Msg.Creator)
	// find subaccount
	subAccAddr, exists := k.keeper.GetSubAccountByOwner(ctx, subAccOwner)
	if !exists {
		return nil, status.Error(codes.NotFound, "subaccount not found")
	}

	bet, err := k.keeper.betKeeper.PrepareBetObject(ctx, msg.Msg.Creator, msg.Msg.Props)
	if err != nil {
		return nil, err
	}

	// here we swap the original creator with the subaccount address
	bet.Creator = subAccAddr.String()

	// make subaccount balance adjustments
	balance, exists := k.keeper.GetBalance(ctx, subAccAddr)
	if !exists {
		panic("state corruption: subaccount balance not found")
	}

	err = balance.Spend(bet.Amount)
	if err != nil {
		return nil, err
	}

	if err := k.keeper.betKeeper.Wager(ctx, bet); err != nil {
		return nil, sdkerrors.Wrapf(bettypes.ErrInWager, "%s", err)
	}

	k.keeper.SetBalance(ctx, subAccAddr, balance)

	msg.EmitEvent(&ctx, subAccOwner.String())

	return &types.MsgWagerResponse{
		Response: &bettypes.MsgWagerResponse{Props: msg.Msg.Props},
	}, nil
}
