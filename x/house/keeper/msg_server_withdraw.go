package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/house/types"
)

// Withdraw performs withdrawal of unused tokens corresponding to a deposit.
func (k msgServer) Withdraw(goCtx context.Context,
	msg *types.MsgWithdraw,
) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depositorAddr, isOnBehalf, err := k.Keeper.ParseWithdrawTicketAndValidate(goCtx, ctx, msg, true)
	if err != nil {
		return nil, err
	}

	id, err := k.Keeper.CalcAndWithdraw(ctx, msg, depositorAddr, isOnBehalf)
	if err != nil {
		return nil, err
	}

	msg.EmitEvent(&ctx, depositorAddr, id)

	return &types.MsgWithdrawResponse{
		ID:                 id,
		MarketUID:          msg.MarketUID,
		ParticipationIndex: msg.ParticipationIndex,
	}, nil
}

// ParseWithdrawTicketAndValidate parses the withdraw payload ticket and validate.
func (k Keeper) ParseWithdrawTicketAndValidate(
	goCtx context.Context,
	ctx sdk.Context,
	msg *types.MsgWithdraw,
	authzAllowed bool,
) (string, bool, error) {
	var payload types.WithdrawTicketPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return "", false, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	isOnBehalf := false
	depositorAddr := msg.Creator
	if payload.DepositorAddress != "" {
		if !authzAllowed {
			return "", false, types.ErrAuthorizationNotAllowed
		}
		depositorAddr = payload.DepositorAddress
		isOnBehalf = true
	}

	if err := payload.Validate(depositorAddr); err != nil {
		return "", false, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	return depositorAddr, isOnBehalf, nil
}

// CalcAndWithdraw calculates the withdrawable amount and withdraws the deposit.
func (k Keeper) CalcAndWithdraw(
	ctx sdk.Context,
	msg *types.MsgWithdraw,
	depositorAddr string,
	isOnBehalf bool,
) (uint64, error) {
	// Get the deposit object
	deposit, found := k.GetDeposit(ctx, depositorAddr, msg.MarketUID, msg.ParticipationIndex)
	if !found {
		return 0, sdkerrors.Wrapf(types.ErrDepositNotFound, ": %s, %d", msg.MarketUID, msg.ParticipationIndex)
	}

	var err error
	msg.Amount, err = k.orderbookKeeper.CalcWithdrawalAmount(ctx,
		depositorAddr,
		msg.MarketUID,
		msg.ParticipationIndex,
		msg.Mode,
		deposit.TotalWithdrawalAmount,
		msg.Amount,
	)
	if err != nil {
		return 0, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if isOnBehalf {
		if err := k.ValidateMsgAuthorization(ctx, msg.Creator, depositorAddr, msg); err != nil {
			return 0, err
		}
	}

	id, err := k.Withdraw(ctx, deposit, msg.Creator, depositorAddr, msg.MarketUID,
		msg.ParticipationIndex, msg.Mode, msg.Amount)
	if err != nil {
		return 0, sdkerrors.Wrap(err, "process withdrawal")
	}

	return id, nil
}
