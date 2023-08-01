package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) CreateSubAccount(
	ctx context.Context,
	request *types.MsgCreateSubAccount,
) (*types.MsgCreateSubAccountResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	err := request.ValidateBasic()
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	moneyToSend, err := sumBalanceUnlocks(sdkContext, request.LockedBalances)
	if err != nil {
		return nil, err
	}

	senderAccount := sdk.MustAccAddressFromBech32(request.Sender)
	subaccountOwner := sdk.MustAccAddressFromBech32(request.SubAccountOwner)
	if _, exists := m.keeper.GetSubAccountByOwner(sdkContext, subaccountOwner); exists {
		return nil, types.ErrSubaccountAlreadyExist
	}

	subaccountID := m.keeper.NextID(sdkContext)

	// ALERT: If someone frontruns the account creation, will be overwritten here
	subaccountAddress := types.NewAddressFromSubaccount(subaccountID)
	subaccountAccount := m.accountKeeper.NewAccountWithAddress(sdkContext, subaccountAddress)
	m.accountKeeper.SetAccount(sdkContext, subaccountAccount)

	err = m.sendCoinsToSubaccount(sdkContext, senderAccount, subaccountAddress, moneyToSend)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send coins")
	}

	m.keeper.SetSubAccountOwner(sdkContext, subaccountAddress, subaccountOwner)
	m.keeper.SetLockedBalances(sdkContext, subaccountAddress, request.LockedBalances)
	m.keeper.SetBalance(sdkContext, subaccountAddress, types.Balance{
		DepositedAmount: moneyToSend,
		SpentAmount:     sdk.ZeroInt(),
		WithdrawmAmount: sdk.ZeroInt(),
		LostAmount:      sdk.ZeroInt(),
	})

	return &types.MsgCreateSubAccountResponse{}, nil
}
