package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/bet/types"
)

func (k msgServer) Wager(
	goCtx context.Context,
	msg *types.MsgWager,
) (*types.MsgWagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	bet, odsMap, err := k.PrepareBetObject(ctx, msg.Creator, msg.Props)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.Wager(ctx, bet, odsMap); err != nil {
		return nil, err
	}

	msg.EmitEvent(&ctx)

	return &types.MsgWagerResponse{Props: msg.Props}, nil
}
