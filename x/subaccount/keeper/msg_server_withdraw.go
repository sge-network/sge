package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) WithdrawUnlockedBalances(ctx context.Context, balances *types.MsgWithdrawUnlockedBalances) (*types.MsgWithdrawUnlockedBalancesResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	sender := sdk.MustAccAddressFromBech32(balances.Sender)
	subAccountAddress, exists := m.keeper.GetSubAccountByOwner(sdkContext, sender)
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	params := m.keeper.GetParams(sdkContext)

	balance, unlockedBalance, bankBalance := m.keeper.getBalances(sdkContext, subAccountAddress, params)

	// calculate withdrawable balance, which is the minimum between the available balance, and
	// what has been unlocked so far. Also, it cannot be greater than the bank balance.
	// Available reports the deposited amount - spent amount - lost amount - withdrawn amount.
	withdrawableBalance := sdk.MinInt(sdk.MinInt(balance.Available(), unlockedBalance), bankBalance.Amount)
	if withdrawableBalance.IsZero() {
		return nil, types.ErrNothingToWithdraw
	}

	balance.WithdrawmAmount = balance.WithdrawmAmount.Add(withdrawableBalance)
	m.keeper.SetBalance(sdkContext, subAccountAddress, balance)

	err := m.bankKeeper.SendCoins(sdkContext, subAccountAddress, sender, sdk.NewCoins(sdk.NewCoin(params.LockedBalanceDenom, withdrawableBalance)))
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawUnlockedBalancesResponse{}, nil
}
