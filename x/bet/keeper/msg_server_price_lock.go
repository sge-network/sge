package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/bet/types"
)

func (k msgServer) PriceLockPoolTopUp(
	goCtx context.Context,
	msg *types.MsgPriceLockPoolTopUp,
) (*types.MsgPriceLockPoolTopUpResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Creator != msg.Funder {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, msg.Funder, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	// if err := k.Keeper.TopUpPriceLockPool(ctx, msg.Funder, msg.Amount); err != nil {
	// 	return nil, err
	// }

	msg.EmitEvent(&ctx)

	return &types.MsgPriceLockPoolTopUpResponse{}, nil
}
