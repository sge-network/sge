package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/sportevent/types"
)

// UpdateSportEvent accepts ticket containing multiple update events and return batch response after processing
func (k msgServer) UpdateSportEvent(goCtx context.Context, msg *types.MsgUpdateSportEventRequest) (*types.MsgUpdateSportEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var updatePayload types.SportEventUpdateTicketPayload
	if err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &updatePayload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	currentData, found := k.Keeper.GetSportEvent(ctx, updatePayload.GetUID())
	if !found {
		return nil, types.ErrEventNotFound
	}

	// if update event is not valid so it is failed
	if err := k.validateEventUpdate(ctx, updatePayload, currentData); err != nil {
		return nil, sdkerrors.Wrap(err, "validate update data")
	}

	// if the status is still pending so it is failed
	if currentData.Status != types.SportEventStatus_SPORT_EVENT_STATUS_UNSPECIFIED {
		return nil, types.ErrCanNotBeAltered
	}

	// replace current data with payload values
	currentData.StartTS = updatePayload.StartTS
	currentData.EndTS = updatePayload.EndTS
	currentData.BetConstraints = &types.EventBetConstraints{
		MinAmount: updatePayload.BetConstraints.MinAmount,
		BetFee:    updatePayload.BetConstraints.BetFee,
	}
	currentData.Active = updatePayload.Active

	// the update event is successful so update the module state
	k.Keeper.SetSportEvent(ctx, currentData)

	response := &types.MsgUpdateSportEventResponse{
		Data: &currentData,
	}
	emitTransactionEvent(ctx, types.TypeMsgUpdateSportEvents, response.Data.UID, msg.Creator)
	return response, nil
}

// validateEventUpdate validates individual event acceptability
func (k msgServer) validateEventUpdate(ctx sdk.Context, updatePayload types.SportEventUpdateTicketPayload, currentData types.SportEvent) error {
	if updatePayload.EndTS <= uint64(ctx.BlockTime().Unix()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid end timestamp for the sport-event")
	}

	if updatePayload.StartTS >= updatePayload.EndTS || updatePayload.StartTS == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid start timestamp for the sport-event")
	}

	params := k.GetParams(ctx)

	// if no bet constraints included in payload, set it as stored before
	if updatePayload.BetConstraints == nil {
		updatePayload.BetConstraints = currentData.BetConstraints
		return nil
	}

	// set current bet fee if the update payload is nil
	if updatePayload.BetConstraints.BetFee.IsNil() {
		updatePayload.BetConstraints.BetFee = currentData.BetConstraints.BetFee
	}

	// set current minimum amount if the update payload is nil
	if updatePayload.BetConstraints.MinAmount.IsNil() {
		updatePayload.BetConstraints.MinAmount = currentData.BetConstraints.MinAmount
	}

	// check the validity constraints as there is no GT method on coin type
	if !(updatePayload.BetConstraints.BetFee.IsLT(params.EventMinBetFee) || updatePayload.BetConstraints.BetFee.IsEqual(params.EventMinBetFee)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event bet fee is out of threshold limit")
	}

	if updatePayload.BetConstraints.MinAmount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min amount can not be negative")
	}
	if updatePayload.BetConstraints.MinAmount.LT(params.EventMinBetAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min bet amount is less than threshold")
	}

	return nil
}
