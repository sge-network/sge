package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/house/types"
)

// Withdraw defines a method for performing a withdrawal of coins of unused amount corresponding to a deposit.
func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id, err := k.Keeper.Withdraw(ctx, msg.Creator, msg.SportEventUID, msg.ParticipationIndex, msg.Mode, msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "process withdrawal")
	}

	return &types.MsgWithdrawResponse{
		ID:                 id,
		SportEventUID:      msg.SportEventUID,
		ParticipationIndex: msg.ParticipationIndex,
	}, nil
}
