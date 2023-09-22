package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/subaccount/types"
)

// Create creates a sub account according to the input message data.
func (k msgServer) Create(
	goCtx context.Context,
	msg *types.MsgCreate,
) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid request")
	}

	moneyToSend, err := sumBalanceUnlocks(ctx, msg.LockedBalances)
	if err != nil {
		return nil, err
	}

	creatorAddr := sdk.MustAccAddressFromBech32(msg.Creator)
	subAccAddr := sdk.MustAccAddressFromBech32(msg.SubAccountOwner)
	if _, exists := k.keeper.GetSubAccountByOwner(ctx, subAccAddr); exists {
		return nil, types.ErrSubaccountAlreadyExist
	}

	subaccountID := k.keeper.NextID(ctx)

	// ALERT: If someone frontruns the account creation, will be overwritten here
	subaccountAddress := types.NewAddressFromSubaccount(subaccountID)
	subaccountAccount := k.keeper.accountKeeper.NewAccountWithAddress(ctx, subaccountAddress)
	k.keeper.accountKeeper.SetAccount(ctx, subaccountAccount)

	err = k.keeper.sendCoinsToSubaccount(ctx, creatorAddr, subaccountAddress, moneyToSend)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to send coins")
	}

	k.keeper.SetSubAccountOwner(ctx, subaccountAddress, subAccAddr)
	k.keeper.SetLockedBalances(ctx, subaccountAddress, msg.LockedBalances)
	k.keeper.SetBalance(ctx, subaccountAddress, types.Balance{
		DepositedAmount: moneyToSend,
		SpentAmount:     sdk.ZeroInt(),
		WithdrawmAmount: sdk.ZeroInt(),
		LostAmount:      sdk.ZeroInt(),
	})

	msg.EmitEvent(&ctx, subaccountAddress.String())

	return &types.MsgCreateResponse{}, nil
}
