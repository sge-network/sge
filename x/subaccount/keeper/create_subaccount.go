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
	if m.keeper.HasSubAccount(sdkContext, subaccountOwner) {
		return nil, types.ErrSubaccountAlreadyExist
	}

	subaccountID := m.keeper.NextID(sdkContext)

	// ALERT: If someone frontruns the account creation, will be overwritten here
	subaccountAddress := types.NewAddressFromSubaccount(subaccountID)
	address := m.accountKeeper.NewAccountWithAddress(sdkContext, subaccountAddress)
	m.accountKeeper.SetAccount(sdkContext, address)

	denom := m.keeper.GetParams(sdkContext).LockedBalanceDenom
	err = m.bankKeeper.SendCoins(sdkContext, senderAccount, subaccountAddress, sdk.NewCoins(sdk.NewCoin(denom, moneyToSend)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to send coins")
	}

	m.keeper.SetSubAccountOwner(sdkContext, subaccountID, subaccountOwner)
	m.keeper.SetLockedBalances(sdkContext, subaccountAddress, request.LockedBalances)
	m.keeper.SetBalance(sdkContext, subaccountID, types.Balance{
		DepositedAmount: moneyToSend,
		SpentAmount:     sdk.ZeroInt(),
		WithdrawmAmount: sdk.ZeroInt(),
		LostAmount:      sdk.ZeroInt(),
	})

	return &types.MsgCreateSubAccountResponse{}, nil
}
