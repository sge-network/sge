package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/rewards/types"
)

func (k msgServer) RewardUser(goCtx context.Context, msg *types.MsgRewardUser) (*types.MsgRewardUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgRewardUserResponse{}, nil
}
