package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) WithdrawUnlockedBalances(ctx context.Context, balances *types.MsgWithdrawUnlockedBalances) (*types.MsgWithdrawUnlockedBalancesResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	sender := sdk.MustAccAddressFromBech32(balances.Sender)
	if !m.keeper.HasSubAccount(sdkContext, sender) {
		return nil, types.ErrSubaccountDoesNotExist
	}

	params := m.keeper.GetParams(sdkContext)
	subaccountID := m.keeper.GetSubAccountByOwner(sdkContext, sender)
	subaccountAddr := types.NewAddressFromSubaccount(subaccountID)

	balance, unlockedBalance, bankBalance := m.getBalances(sdkContext, subaccountID, subaccountAddr, params)

	// calculate withdrawable balance
	withdrawableBalance := sdk.MinInt(sdk.MinInt(balance.DepositedAmount, unlockedBalance), bankBalance.Amount)

	balance.WithdrawmAmount = balance.WithdrawmAmount.Add(withdrawableBalance)
	m.keeper.SetBalance(sdkContext, subaccountID, balance)

	err := m.bankKeeper.SendCoins(sdkContext, subaccountAddr, sender, sdk.NewCoins(sdk.NewCoin(params.LockedBalanceDenom, withdrawableBalance)))
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawUnlockedBalancesResponse{}, nil
}

// getBalances returns the balance, unlocked balance and bank balance of a subaccount
func (m msgServer) getBalances(sdkContext sdk.Context, subaccountID uint64, subaccountAddr sdk.AccAddress, params types.Params) (types.Balance, sdk.Int, sdk.Coin) {
	balance := m.keeper.GetBalance(sdkContext, subaccountID)
	unlockedBalance := m.keeper.GetUnlockedBalance(sdkContext, subaccountID)
	bankBalance := m.bankKeeper.GetBalance(sdkContext, subaccountAddr, params.LockedBalanceDenom)

	return balance, unlockedBalance, bankBalance
}
