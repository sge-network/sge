package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/house/types"
)

// Deposit performs deposit operation to participate as a house in a specific market/order book
func (k msgServer) Deposit(goCtx context.Context,
	msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateDeposit(ctx, msg); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid deposit")
	}

	participationIndex, err := k.Keeper.Deposit(ctx, msg.Creator, msg.MarketUID, msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	return &types.MsgDepositResponse{
		MarketUID:          msg.MarketUID,
		ParticipationIndex: participationIndex,
	}, nil
}
