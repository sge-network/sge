package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) CreateSubAccount(
	ctx context.Context,
	request *types.MsgCreateSubAccountRequest,
) (*types.MsgCreateAccountResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	moneyToSend := sdk.NewInt(0)

	err := request.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	for _, balanceUnlock := range request.LockedBalances {
		if balanceUnlock.UnlockTime.Unix() < sdkContext.BlockTime().Unix() {
			return nil, types.ErrUnlockTokenTimeExpired
		}

		moneyToSend = moneyToSend.Add(balanceUnlock.Amount)
	}

	senderAccount, _ := sdk.AccAddressFromBech32(request.Sender)
	subaccountOwner, _ := sdk.AccAddressFromBech32(request.SubAccountOwner)
	if m.keeper.HasSubAccount(sdkContext, subaccountOwner) {
		return nil, types.ErrSubaccountAlreadyExist
	}

	subaccountID := m.keeper.NextID(sdkContext)

	subaccountAddress := types.NewAddressFromSubaccount(subaccountID)
	address := m.accountKeeper.NewAccountWithAddress(sdkContext, subaccountAddress)
	m.accountKeeper.SetAccount(sdkContext, address)

	err = m.bankKeeper.SendCoins(sdkContext, senderAccount, subaccountAddress, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, moneyToSend)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to send coins")
	}

	m.keeper.SetSubAccountOwner(sdkContext, subaccountID, subaccountOwner)
	m.keeper.SetLockedBalances(sdkContext, subaccountAddress, request.LockedBalances)
	m.keeper.SetBalance(sdkContext, subaccountAddress, types.Balance{
		DepositedAmount: moneyToSend,
		SpentAmount:     sdk.ZeroInt(),
		WithdrawmAmount: sdk.ZeroInt(),
		LostAmount:      sdk.ZeroInt(),
	})

	return &types.MsgCreateAccountResponse{}, nil
}
