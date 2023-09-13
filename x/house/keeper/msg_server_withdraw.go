package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/house/types"
)

// Withdraw performs withdrawal of unused tokens corresponding to a deposit.
func (k msgServer) Withdraw(goCtx context.Context,
	msg *types.MsgWithdraw,
) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.WithdrawTicketPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	isOnBehalf := false
	depositorAddr := msg.Creator
	if payload.DepositorAddress != "" {
		depositorAddr = payload.DepositorAddress
		isOnBehalf = true
	}

	if err := payload.Validate(depositorAddr); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	// Get the deposit object
	deposit, found := k.GetDeposit(ctx, depositorAddr, msg.MarketUID, msg.ParticipationIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDepositNotFound, ": %s, %d", msg.MarketUID, msg.ParticipationIndex)
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
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if isOnBehalf {
		if err := utils.ValidateMsgAuthorization(k.authzKeeper, ctx, msg.Creator, payload.DepositorAddress, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return nil, err
		}
	}

	id, err := k.Keeper.Withdraw(ctx, deposit, msg.Creator, depositorAddr, msg.MarketUID,
		msg.ParticipationIndex, msg.Mode, msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "process withdrawal")
	}

	msg.EmitEvent(&ctx, depositorAddr, id)

	return &types.MsgWithdrawResponse{
		ID:                 id,
		MarketUID:          msg.MarketUID,
		ParticipationIndex: msg.ParticipationIndex,
	}, nil
}
