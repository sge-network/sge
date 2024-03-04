package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/subaccount/types"
)

// TopUp increases the balance of sub account according to the input data.
func (k msgServer) TopUp(goCtx context.Context, msg *types.MsgTopUp) (*types.MsgTopUpResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	subAccAddr, err := k.keeper.TopUp(ctx, msg.Creator, msg.Address, msg.LockedBalances)
	if err != nil {
		return nil, err
	}

	msg.EmitEvent(&ctx, subAccAddr)

	return &types.MsgTopUpResponse{}, nil
}

// WithdrawUnlockedBalances withdraws the unlocked balance of sub account according to the input account.
func (k msgServer) WithdrawUnlockedBalances(goCtx context.Context, msg *types.MsgWithdrawUnlockedBalances) (*types.MsgWithdrawUnlockedBalancesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ownerAddr := sdk.MustAccAddressFromBech32(msg.Creator)
	subAccAddr, exists := k.keeper.GetSubAccountByOwner(ctx, ownerAddr)
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	err := k.keeper.withdrawUnlocked(ctx, subAccAddr, ownerAddr)
	if err != nil {
		return nil, err
	}

	msg.EmitEvent(&ctx, subAccAddr.String())

	return &types.MsgWithdrawUnlockedBalancesResponse{}, nil
}
