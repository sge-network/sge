package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/sportevent/types"
)

// UpdateSportEvent accepts ticket containing multiple update events and return batch response after processing
func (k msgServer) UpdateSportEvent(goCtx context.Context, msg *types.MsgUpdateSportEvent) (*types.MsgUpdateSportEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var updatePayload types.SportEventUpdateTicketPayload
	if err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &updatePayload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	currentData, found := k.Keeper.GetSportEvent(ctx, updatePayload.GetUID())
	if !found {
		return nil, types.ErrEventNotFound
	}

	// if current data is not active or inactive it is not updatable
	// active status can be changed to inactive or vice versa in the updating
	if currentData.Status != types.SportEventStatus_SPORT_EVENT_STATUS_ACTIVE &&
		currentData.Status != types.SportEventStatus_SPORT_EVENT_STATUS_INACTIVE {
		return nil, types.ErrCanNotBeAltered
	}

	// if update event is not valid so it is failed
	params := k.GetParams(ctx)

	if err := updatePayload.Validate(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "validate update data")
	}

	// replace current data with payload values
	currentData.StartTS = updatePayload.StartTS
	currentData.EndTS = updatePayload.EndTS
	currentData.BetConstraints = params.NewEventBetConstraints(updatePayload.MinBetAmount, updatePayload.BetFee)
	currentData.Status = updatePayload.Status

	// the update event is successful so update the module state
	k.Keeper.SetSportEvent(ctx, currentData)

	response := &types.MsgUpdateSportEventResponse{
		Data: &currentData,
	}
	emitTransactionEvent(ctx, types.TypeMsgUpdateSportEvents, response.Data.UID, response.Data.BookID, msg.Creator)
	return response, nil
}
