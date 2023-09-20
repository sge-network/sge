package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/subaccount/types"
)

// TopUp increases the balance of sub account according to the input data.
func (m msgServer) TopUp(ctx context.Context, msg *types.MsgTopUp) (*types.MsgTopUpResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	moneyToAdd, err := sumBalanceUnlocks(sdkContext, msg.LockedBalances)
	if err != nil {
		return nil, err
	}

	creatorAddr := sdk.MustAccAddressFromBech32(msg.Creator)
	subaccountOwner := sdk.MustAccAddressFromBech32(msg.SubAccount)

	subAccAddress, exists := m.keeper.GetSubAccountByOwner(sdkContext, subaccountOwner)
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}
	balance, exists := m.keeper.GetBalance(sdkContext, subAccAddress)
	if !exists {
		panic("data corruption: subaccount exists but balance does not")
	}

	balance.DepositedAmount = balance.DepositedAmount.Add(moneyToAdd)
	m.keeper.SetBalance(sdkContext, subAccAddress, balance)
	m.keeper.SetLockedBalances(sdkContext, subAccAddress, msg.LockedBalances)

	err = m.keeper.sendCoinsToSubaccount(sdkContext, creatorAddr, subAccAddress, moneyToAdd)
	if err != nil {
		return nil, fmt.Errorf("unable to send coins: %w", err)
	}

	return &types.MsgTopUpResponse{}, nil
}

// WithdrawUnlockedBalances withdraws the unlocked balance of sub account according to the input account.
func (m msgServer) WithdrawUnlockedBalances(ctx context.Context, balances *types.MsgWithdrawUnlockedBalances) (*types.MsgWithdrawUnlockedBalancesResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	creatorAddr := sdk.MustAccAddressFromBech32(balances.Creator)
	subAccountAddress, exists := m.keeper.GetSubAccountByOwner(sdkContext, creatorAddr)
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	balance, unlockedBalance, bankBalance := m.keeper.getBalances(sdkContext, subAccountAddress)

	// calculate withdrawable balance, which is the minimum between the available balance, and
	// what has been unlocked so far. Also, it cannot be greater than the bank balance.
	// Available reports the deposited amount - spent amount - lost amount - withdrawn amount.
	withdrawableBalance := sdk.MinInt(sdk.MinInt(balance.Available(), unlockedBalance), bankBalance.Amount)
	if withdrawableBalance.IsZero() {
		return nil, types.ErrNothingToWithdraw
	}

	balance.WithdrawmAmount = balance.WithdrawmAmount.Add(withdrawableBalance)
	m.keeper.SetBalance(sdkContext, subAccountAddress, balance)

	err := m.keeper.bankKeeper.SendCoins(sdkContext, subAccountAddress, creatorAddr, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, withdrawableBalance)))
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawUnlockedBalancesResponse{}, nil
}
