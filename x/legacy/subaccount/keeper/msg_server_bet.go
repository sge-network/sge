package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/app/params"

	bettypes "github.com/sge-network/sge/x/legacy/bet/types"
	"github.com/sge-network/sge/x/legacy/subaccount/types"
)

func (k msgServer) Wager(goCtx context.Context, msg *types.MsgWager) (*types.MsgWagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.Keeper.GetWagerEnabled(ctx) {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "currently the subacount wager tx is not enabled")
	}

	subAccOwner := sdk.MustAccAddressFromBech32(msg.Creator)
	// find subaccount
	subAccAddr, exists := k.Keeper.GetSubaccountByOwner(ctx, subAccOwner)
	if !exists {
		return nil, status.Error(codes.NotFound, "subaccount not found")
	}

	payload := &types.SubAccWagerTicketPayload{}
	err := k.Keeper.ovmKeeper.VerifyTicketUnmarshal(ctx, msg.Ticket, &payload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if msg.Creator != payload.Msg.Creator {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "message creator should be the same as the sub message creator%s", msg.Creator)
	}

	bet, oddsMap, err := k.Keeper.betKeeper.PrepareBetObject(ctx, payload.Msg.Creator, payload.Msg.Props)
	if err != nil {
		return nil, err
	}

	if err := payload.Validate(bet.Amount); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	mainAccBalance := k.Keeper.bankKeeper.GetBalance(
		ctx,
		sdk.MustAccAddressFromBech32(bet.Creator),
		params.DefaultBondDenom)
	if mainAccBalance.Amount.LT(payload.MainaccDeductAmount) {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "not enough balance in main account")
	}

	if err := k.Keeper.withdrawLockedAndUnlocked(ctx, subAccAddr, subAccOwner, payload.SubaccDeductAmount); err != nil {
		return nil, err
	}

	if err := k.Keeper.betKeeper.Wager(ctx, bet, oddsMap); err != nil {
		return nil, err
	}

	msg.EmitEvent(&ctx, payload.Msg, subAccOwner.String())

	return &types.MsgWagerResponse{
		Response: &bettypes.MsgWagerResponse{Props: payload.Msg.Props},
	}, nil
}
