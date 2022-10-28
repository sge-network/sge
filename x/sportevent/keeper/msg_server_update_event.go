package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/sportevent/types"
)

// UpdateEvent accepts ticket containing multiple update events and return batch response after processing
func (k msgServer) UpdateEvent(goCtx context.Context, msg *types.MsgUpdateEvent) (*types.MsgSportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var updateEvents types.SportEventUpdateTicket
	if err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &updateEvents); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	success := make([]string, 0, len(updateEvents.Events))
	failed := make([]*types.FailedEvent, 0)
	for _, event := range updateEvents.Events {
		sportEvent, found := k.Keeper.GetSportEvent(ctx, event.GetUID())

		// if sport event is not found so it is failed
		if !found {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: types.ErrEventNotFound.Error(),
			})
			continue
		}

		// if update event is not valid so it is failed
		if err := k.validateEventUpdate(ctx, event, sportEvent); err != nil {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: err.Error(),
			})
			continue
		}

		// if the status is still pending so it is failed
		if sportEvent.Status != types.SportEventStatus_STATUS_PENDING {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: types.ErrCanNotBeAltered.Error(),
			})
			continue
		}

		// the update event is successful so update the module state
		k.Keeper.SetSportEvent(ctx, types.SportEvent{
			UID:            sportEvent.UID,
			OddsUIDs:       sportEvent.OddsUIDs,
			WinnerOddsUIDs: sportEvent.WinnerOddsUIDs,
			Status:         sportEvent.Status,
			ResolutionTS:   sportEvent.ResolutionTS,
			Creator:        sportEvent.Creator,
			StartTS:        event.StartTS,
			EndTS:          event.EndTS,
			BetConstraints: &types.EventBetConstraints{
				MaxBetCap:          event.BetConstraints.MaxBetCap,
				MinAmount:          event.BetConstraints.MinAmount,
				BetFee:             event.BetConstraints.BetFee,
				CurrentTotalAmount: sportEvent.BetConstraints.CurrentTotalAmount,
			},
			Active: event.Active,
		})

		// update success events list
		success = append(success, sportEvent.UID)
	}

	response := &types.MsgSportResponse{
		SuccessEvents: success,
		FailedEvents:  failed,
	}
	emitTransactionEvent(ctx, types.TypeMsgUpdateSportEvents, response, msg.Creator)
	return response, nil
}

// validateEventUpdate validates individual event acceptability
func (k msgServer) validateEventUpdate(ctx sdk.Context, event, previousEvent types.SportEvent) error {
	if event.EndTS <= uint64(ctx.BlockTime().Unix()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid end timestamp for the sport event")
	}

	if event.StartTS >= event.EndTS || event.StartTS == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid start timestamp for the sport event")
	}

	params := k.GetParams(ctx)

	if event.BetConstraints == nil {
		event.BetConstraints = previousEvent.BetConstraints
		return nil
	}

	//init individual params if any one of them is nil
	if event.BetConstraints.BetFee.IsNil() {
		event.BetConstraints.BetFee = previousEvent.BetConstraints.BetFee
	}
	if event.BetConstraints.MaxBetCap.IsNil() {
		event.BetConstraints.MaxBetCap = previousEvent.BetConstraints.MaxBetCap
	}
	if event.BetConstraints.MinAmount.IsNil() {
		event.BetConstraints.MinAmount = previousEvent.BetConstraints.MinAmount
	}

	//check the validity constraints
	if !(event.BetConstraints.BetFee.IsLT(params.EventMinBetFee) || event.BetConstraints.BetFee.IsEqual(params.EventMinBetFee)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event bet fee is out of threshold limit")
	}

	if event.BetConstraints.MinAmount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min amount can not be negative")
	}
	if event.BetConstraints.MinAmount.LT(params.EventMinBetAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min bet amount is less than threshold")
	}

	if event.BetConstraints.MaxBetCap.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event max bet can not be negative")
	}
	if event.BetConstraints.MaxBetCap.GT(params.EventMaxBetCap) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event max bet cap is greater than threshold")
	}
	if event.BetConstraints.MinAmount.GTE(event.BetConstraints.MaxBetCap) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "min bet amount cannot be greater than or equals to to max bet capacity")
	}
	return nil
}
