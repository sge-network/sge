package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/legacy/house/types"
)

// Deposit performs deposit operation to participate as a house in a specific market/order book
func (k msgServer) Deposit(goCtx context.Context,
	msg *types.MsgDeposit,
) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depositorAddr, err := k.ParseDepositTicketAndValidate(goCtx, ctx, msg, true)
	if err != nil {
		return nil, err
	}

	participationIndex, feeAmount, err := k.Keeper.Deposit(
		ctx,
		msg.Creator,
		depositorAddr,
		msg.MarketUID,
		msg.Amount,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	msg.EmitEvent(&ctx, depositorAddr, participationIndex, feeAmount)

	return &types.MsgDepositResponse{
		MarketUID:          msg.MarketUID,
		ParticipationIndex: participationIndex,
	}, nil
}

// ParseDepositTicketAndValidate parses the deposit payload ticket and validate.
func (k Keeper) ParseDepositTicketAndValidate(
	goCtx context.Context,
	ctx sdk.Context,
	msg *types.MsgDeposit,
	authzAllowed bool,
) (string, error) {
	params := k.GetParams(ctx)
	if err := msg.ValidateSanity(ctx, &params); err != nil {
		return "", sdkerrors.Wrap(err, "invalid deposit")
	}

	var payload types.DepositTicketPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return "", sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	depositorAddr := msg.Creator
	if payload.DepositorAddress != "" &&
		payload.DepositorAddress != msg.Creator {
		if !authzAllowed {
			return "", types.ErrAuthorizationNotAllowed
		}
		if err := utils.ValidateMsgAuthorization(ctx, k.authzKeeper, msg.Creator, payload.DepositorAddress, msg,
			types.ErrAuthorizationNotFound, types.ErrAuthorizationNotAccepted); err != nil {
			return "", err
		}
		depositorAddr = payload.DepositorAddress
	}

	if err := payload.Validate(depositorAddr); err != nil {
		return "", sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	return depositorAddr, nil
}
