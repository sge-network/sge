package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
)

func (k msgServer) Wager(
	goCtx context.Context,
	msg *types.MsgWager,
) (*types.MsgWagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	bet, err := k.PrepareBetObject(ctx, msg.Creator, msg.Props)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.Wager(ctx, bet); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInWager, "%s", err)
	}

	msg.EmitEvent(&ctx)

	return &types.MsgWagerResponse{Props: msg.Props}, nil
}
