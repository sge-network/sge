package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/subaccount/types"
)

// Create creates a sub account according to the input message data.
func (k msgServer) Create(
	ctx context.Context,
	request *types.MsgCreate,
) (*types.MsgCreateResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	err := request.ValidateBasic()
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid request")
	}

	moneyToSend, err := sumBalanceUnlocks(sdkContext, request.LockedBalances)
	if err != nil {
		return nil, err
	}

	creatorAddr := sdk.MustAccAddressFromBech32(request.Creator)
	subaccountOwner := sdk.MustAccAddressFromBech32(request.SubAccountOwner)
	if _, exists := k.keeper.GetSubAccountByOwner(sdkContext, subaccountOwner); exists {
		return nil, types.ErrSubaccountAlreadyExist
	}

	subaccountID := k.keeper.NextID(sdkContext)

	// ALERT: If someone frontruns the account creation, will be overwritten here
	subaccountAddress := types.NewAddressFromSubaccount(subaccountID)
	subaccountAccount := k.keeper.accountKeeper.NewAccountWithAddress(sdkContext, subaccountAddress)
	k.keeper.accountKeeper.SetAccount(sdkContext, subaccountAccount)

	err = k.keeper.sendCoinsToSubaccount(sdkContext, creatorAddr, subaccountAddress, moneyToSend)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to send coins")
	}

	k.keeper.SetSubAccountOwner(sdkContext, subaccountAddress, subaccountOwner)
	k.keeper.SetLockedBalances(sdkContext, subaccountAddress, request.LockedBalances)
	k.keeper.SetBalance(sdkContext, subaccountAddress, types.Balance{
		DepositedAmount: moneyToSend,
		SpentAmount:     sdk.ZeroInt(),
		WithdrawmAmount: sdk.ZeroInt(),
		LostAmount:      sdk.ZeroInt(),
	})

	return &types.MsgCreateResponse{}, nil
}
