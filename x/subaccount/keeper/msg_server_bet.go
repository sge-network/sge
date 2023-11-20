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

	subAccOwner := sdk.MustAccAddressFromBech32(msg.Creator)
	// find subaccount
	subAccAddr, exists := k.keeper.GetSubAccountByOwner(ctx, subAccOwner)
	if !exists {
		return nil, status.Error(codes.NotFound, "subaccount not found")
	}

	payload := &types.SubAccWagerTicketPayload{}
	err := k.keeper.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Ticket, &payload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if msg.Creator != payload.Msg.Creator {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "message creator should be the same as the sub message creator%s", msg.Creator)
	}

	bet, oddsMap, err := k.keeper.betKeeper.PrepareBetObject(ctx, payload.Msg.Creator, payload.Msg.Props)
	if err != nil {
		return nil, err
	}

	if err := payload.Validate(bet.Amount); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	mainAccBalance := k.keeper.bankKeeper.GetBalance(
		ctx,
		sdk.MustAccAddressFromBech32(bet.Creator),
		params.DefaultBondDenom)
	if mainAccBalance.Amount.LT(payload.MainaccDeductAmount) {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "not enough balance in main account")
	}

	accSummary, unlockedBalance, _ := k.keeper.getAccountSummary(ctx, subAccAddr)
	if unlockedBalance.GTE(payload.SubaccDeductAmount) {
		if err := k.keeper.bankKeeper.SendCoins(ctx,
			subAccAddr,
			sdk.MustAccAddressFromBech32(msg.Creator),
			sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, payload.SubaccDeductAmount))); err != nil {
			return nil, sdkerrors.Wrapf(types.ErrSendCoinError, "error sending coin from subaccount to main account %s", err)
		}
	} else {
		lockedAmountToWithdraw := payload.SubaccDeductAmount.Sub(unlockedBalance)

		if err := accSummary.Withdraw(lockedAmountToWithdraw); err != nil {
			return nil, sdkerrors.Wrapf(types.ErrWithdrawLocked, "%s", err)
		}

		if err := k.keeper.bankKeeper.SendCoins(ctx, subAccAddr, subAccOwner,
			sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, lockedAmountToWithdraw))); err != nil {
			return nil, sdkerrors.Wrapf(types.ErrSendCoinError, "error sending coin from subaccount to main account %s", err)
		}

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

	if err := k.keeper.betKeeper.Wager(ctx, bet, oddsMap); err != nil {
		return nil, err
	}

	k.keeper.SetAccountSummary(ctx, subAccAddr, accSummary)

	msg.EmitEvent(&ctx, payload.Msg, subAccOwner.String())

	return &types.MsgWagerResponse{
		Response: &bettypes.MsgWagerResponse{Props: payload.Msg.Props},
	}, nil
}
