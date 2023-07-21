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
	if !m.keeper.HasSubAccount(sdkContext, subaccountOwner) {
		return nil, types.ErrSubaccountDoesNotExist
	}

	subaccountID := m.keeper.GetSubAccountByOwner(sdkContext, subaccountOwner)
	balance := m.keeper.GetBalance(sdkContext, subaccountID)

	balance.DepositedAmount = balance.DepositedAmount.Add(moneyToAdd)
	m.keeper.SetBalance(sdkContext, subaccountID, balance)
	m.keeper.SetLockedBalances(sdkContext, subaccountID, msg.LockedBalances)

	err = m.sendCoinsToSubaccount(sdkContext, sender, subaccountID, moneyToAdd)
	if err != nil {
		return nil, fmt.Errorf("unable to send coins: %w", err)
	}

	return &types.MsgTopUpResponse{}, nil
}
