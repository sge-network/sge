package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	housetypes "github.com/sge-network/sge/x/house/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) HouseDeposit(goCtx context.Context, msg *types.MsgHouseDeposit) (*types.MsgHouseDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if subaccount exists
	subAccountAddr, exists := m.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(msg.Msg.Creator))
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	if err := m.houseDeposit(ctx, msg.Msg); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	// send house deposit from subaccount on behalf of the owner
	participationIndex, err := m.keeper.houseKeeper.Deposit(
		ctx,
		subAccountAddr.String(),
		subAccountAddr.String(),
		msg.Msg.MarketUID,
		msg.Msg.Amount,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}
	msg.Msg.EmitEvent(&ctx, subAccountAddr.String(), participationIndex)
	return &types.MsgHouseDepositResponse{
		Response: &housetypes.MsgDepositResponse{
			MarketUID:          msg.Msg.MarketUID,
			ParticipationIndex: participationIndex,
		},
	}, nil
}

// TODO: This is a copy of the Deposit function from x/house/keeper/msg_server_deposit.go
func (m msgServer) houseDeposit(ctx sdk.Context, msg *housetypes.MsgDeposit) error {

	params := m.keeper.houseKeeper.GetParams(ctx)
	if err := msg.ValidateSanity(ctx, &params); err != nil {
		return sdkerrors.Wrap(err, "invalid deposit")
	}

	var payload housetypes.DepositTicketPayload
	if err := m.keeper.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Ticket, &payload); err != nil {
		return sdkerrors.Wrapf(housetypes.ErrInTicketVerification, "%s", err)
	}

	depositorAddr := msg.Creator

	if err := payload.Validate(depositorAddr); err != nil {
		return sdkerrors.Wrapf(housetypes.ErrInTicketPayloadValidation, "%s", err)
	}

	return nil
}

func (m msgServer) HouseWithdraw(ctx context.Context, withdraw *types.MsgHouseWithdraw) (*types.MsgHouseWithdrawResponse, error) {
	// TODO implement me
	panic("implement me")
}
