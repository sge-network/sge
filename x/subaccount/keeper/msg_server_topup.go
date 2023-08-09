package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) TopUp(ctx context.Context, msg *types.MsgTopUp) (*types.MsgTopUpResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	moneyToAdd, err := sumBalanceUnlocks(sdkContext, msg.LockedBalances)
	if err != nil {
		return nil, err
	}

	sender := sdk.MustAccAddressFromBech32(msg.Sender)
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

	err = m.keeper.sendCoinsToSubaccount(sdkContext, sender, subAccAddress, moneyToAdd)
	if err != nil {
		return nil, fmt.Errorf("unable to send coins: %w", err)
	}

	return &types.MsgTopUpResponse{}, nil
}
