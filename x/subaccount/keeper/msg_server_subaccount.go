package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/subaccount/types"
)

// Create creates a sub account according to the input message data.
func (k msgServer) Create(
	goCtx context.Context,
	msg *types.MsgCreate,
) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid request")
	}

	subAccAddr, err := k.keeper.CreateSubAccount(ctx, msg.Creator, msg.SubAccountOwner, msg.LockedBalances)
	if err != nil {
		return nil, err
	}

	msg.EmitEvent(&ctx, subAccAddr)

	return &types.MsgCreateResponse{}, nil
}
