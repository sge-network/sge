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

	// if update event is not valid so it is failed
	params := k.GetParams(ctx)

	if err := updatePayload.Validate(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "validate update data")
	}

	// if the status is still pending so it is failed
	if currentData.Status != types.SportEventStatus_SPORT_EVENT_STATUS_PENDING {
		return nil, types.ErrCanNotBeAltered
	}

	// replace current data with payload values
	currentData.StartTS = updatePayload.StartTS
	currentData.EndTS = updatePayload.EndTS
	currentData.BetConstraints = updatePayload.GetBetConstraints()
	currentData.Active = updatePayload.Active

	// the update event is successful so update the module state
	k.Keeper.SetSportEvent(ctx, currentData)

	response := &types.MsgUpdateSportEventResponse{
		Data: &currentData,
	}
	emitTransactionEvent(ctx, types.TypeMsgUpdateSportEvents, response.Data.UID, msg.Creator)
	return response, nil
}
