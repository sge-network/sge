package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/app/params"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (k msgServer) Wager(goCtx context.Context, msg *types.MsgWager) (*types.MsgWagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	subAccOwner := sdk.MustAccAddressFromBech32(msg.Msg.Creator)
	// find subaccount
	subAccAddr, exists := k.keeper.GetSubAccountByOwner(ctx, subAccOwner)
	if !exists {
		return nil, status.Error(codes.NotFound, "subaccount not found")
	}

	bet, oddsMap, err := k.keeper.betKeeper.PrepareBetObject(ctx, msg.Msg.Creator, msg.Msg.Props)
	if err != nil {
		return nil, err
	}

	mainAccBalance := k.keeper.bankKeeper.GetBalance(
		ctx,
		sdk.MustAccAddressFromBech32(bet.Creator),
		params.DefaultBondDenom)
	if mainAccBalance.Amount.LT(msg.MainaccDeductAmount) {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "not enough balance in main account")
	}

	accSummary, unlockedBalance, _ := k.keeper.getAccountSummary(ctx, subAccAddr)
	if unlockedBalance.GTE(msg.SubaccDeductAmount) {
		k.keeper.bankKeeper.SendCoins(ctx,
			subAccAddr,
			sdk.MustAccAddressFromBech32(msg.Msg.Creator),
			sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, msg.SubaccDeductAmount)))
	} else {
		lockedAmountToWithdraw := msg.SubaccDeductAmount.Sub(unlockedBalance)
		accSummary.Withdraw(lockedAmountToWithdraw)
		k.keeper.bankKeeper.SendCoins(ctx, subAccAddr, subAccOwner,
			sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, lockedAmountToWithdraw)))

		// calculate locked balances
		lockedBalances := k.keeper.GetLockedBalances(ctx, subAccAddr)
		updatedLockedBalances := []types.LockedBalance{}
		for _, lb := range lockedBalances {
			if lockedAmountToWithdraw.LTE(sdkmath.ZeroInt()) {
				break
			}

			if lb.Amount.GT(lockedAmountToWithdraw) {
				lb.Amount = lb.Amount.Sub(lockedAmountToWithdraw)
				updatedLockedBalances = append(updatedLockedBalances, lb)

				lockedAmountToWithdraw = sdkmath.ZeroInt()
				break
			} else {
				lb.Amount = sdkmath.ZeroInt()
				updatedLockedBalances = append(updatedLockedBalances, lb)

				lockedAmountToWithdraw = lockedAmountToWithdraw.Sub(lb.Amount)
			}

		}

		if lockedAmountToWithdraw.GT(sdkmath.ZeroInt()) {
			return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "not enough balance in sub account")
		}

		k.keeper.SetLockedBalances(ctx, subAccAddr, updatedLockedBalances)
	}

	// if err := accSummary.Spend(msg.SubaccDeductAmount); err != nil {
	// 	return nil, err
	// }

	if err := k.keeper.betKeeper.Wager(ctx, bet, oddsMap); err != nil {
		return nil, err
	}

	k.keeper.SetAccountSummary(ctx, subAccAddr, accSummary)

	msg.EmitEvent(&ctx, subAccOwner.String())

	return &types.MsgWagerResponse{
		Response: &bettypes.MsgWagerResponse{Props: msg.Msg.Props},
	}, nil
}
