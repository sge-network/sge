package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
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

	creatorAddr := sdk.MustAccAddressFromBech32(msg.Creator)
	subAccAddr, exists := k.keeper.GetSubAccountByOwner(ctx, creatorAddr)
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	balance, unlockedBalance, bankBalance := k.keeper.getBalances(ctx, subAccAddr)

	// calculate withdrawable balance, which is the minimum between the available balance, and
	// what has been unlocked so far. Also, it cannot be greater than the bank balance.
	// Available reports the deposited amount - spent amount - lost amount - withdrawn amount.
	withdrawableBalance := sdk.MinInt(sdk.MinInt(balance.Available(), unlockedBalance), bankBalance.Amount)
	if withdrawableBalance.IsZero() {
		return nil, types.ErrNothingToWithdraw
	}

	balance.WithdrawnAmount = balance.WithdrawnAmount.Add(withdrawableBalance)
	k.keeper.SetBalance(ctx, subAccAddr, balance)

	err := k.keeper.bankKeeper.SendCoins(ctx, subAccAddr, creatorAddr, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, withdrawableBalance)))
	if err != nil {
		return nil, err
	}

	msg.EmitEvent(&ctx, subAccAddr.String())

	return &types.MsgWithdrawUnlockedBalancesResponse{}, nil
}
