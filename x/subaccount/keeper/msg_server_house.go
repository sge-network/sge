package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	housetypes "github.com/sge-network/sge/x/house/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (k msgServer) HouseDeposit(goCtx context.Context, msg *types.MsgHouseDeposit) (*types.MsgHouseDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.keeper.GetDepositEnabled(ctx) {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "currently the subacount deposit tx is not enabled")
	}

	// check if subaccount exists
	subAccAddr, exists := k.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(msg.Msg.Creator))
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	// get subaccount balance, and check if it can spend
	balance, exists := k.keeper.GetAccountSummary(ctx, subAccAddr)
	if !exists {
		panic("data corruption: subaccount balance not found")
	}

	// parse the ticket payload and create deposit, setting the authz allowed as false
	_, err := k.keeper.houseKeeper.ParseDepositTicketAndValidate(goCtx, ctx, msg.Msg, false)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	if err := balance.Spend(msg.Msg.Amount); err != nil {
		return nil, err
	}

	// send house deposit from subaccount on behalf of the owner
	participationIndex, feeAmount, err := k.keeper.houseKeeper.Deposit(
		ctx,
		msg.Msg.Creator,
		subAccAddr.String(),
		msg.Msg.MarketUID,
		msg.Msg.Amount,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	// update subaccount balance
	k.keeper.SetAccountSummary(ctx, subAccAddr, balance)

	// emit event
	msg.EmitEvent(&ctx, subAccAddr.String(), participationIndex, feeAmount)

	return &types.MsgHouseDepositResponse{
		Response: &housetypes.MsgDepositResponse{
			MarketUID:          msg.Msg.MarketUID,
			ParticipationIndex: participationIndex,
		},
	}, nil
}

func (k msgServer) HouseWithdraw(goCtx context.Context, msg *types.MsgHouseWithdraw) (*types.MsgHouseWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if subaccount exists
	subAccAddr, exists := k.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(msg.Msg.Creator))
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	subAccountBalance, exists := k.keeper.GetAccountSummary(ctx, subAccAddr)
	if !exists {
		panic("data corruption: subaccount balance not found")
	}

	_, _, err := k.keeper.houseKeeper.ParseWithdrawTicketAndValidate(goCtx, msg.Msg, true)
	if err != nil {
		return nil, err
	}

	id, err := k.keeper.houseKeeper.CalcAndWithdraw(ctx, msg.Msg, subAccAddr.String(), false)
	if err != nil {
		return nil, err
	}

	err = subAccountBalance.Unspend(msg.Msg.Amount)
	if err != nil {
		panic("data corruption: it must be possible to unspend an house withdrawal")
	}

	k.keeper.SetAccountSummary(ctx, subAccAddr, subAccountBalance)

	msg.Msg.EmitEvent(&ctx, subAccAddr.String(), id)

	return &types.MsgHouseWithdrawResponse{
		Response: &housetypes.MsgWithdrawResponse{
			ID:                 id,
			MarketUID:          msg.Msg.MarketUID,
			ParticipationIndex: msg.Msg.ParticipationIndex,
		},
	}, nil
}
