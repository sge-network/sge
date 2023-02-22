package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/sportevent/types"
)

// ResolveSportEvent accepts ticket containing multiple resolution events and return batch response after processing
func (k msgServer) ResolveSportEvent(goCtx context.Context, msg *types.MsgResolveSportEvent) (*types.MsgResolveSportEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var resolutionPayload types.SportEventResolutionTicketPayload
	err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	if err := k.processSportEventResolution(ctx, &resolutionPayload); err != nil {
		return nil, sdkerrors.Wrap(err, "process resolution event")
	}

	sportEvent, _ := k.getSportEventToResolve(ctx, resolutionPayload)
	response := &types.MsgResolveSportEventResponse{
		Data: &sportEvent,
	}
	emitTransactionEvent(ctx, types.TypeMsgResolveSportEvents, response.Data.UID, msg.Creator)
	return response, nil
}

func (k msgServer) processSportEventResolution(ctx sdk.Context, resolutionPayload *types.SportEventResolutionTicketPayload) error {
	if err := resolutionPayload.Validate(); err != nil {
		return sdkerrors.Wrap(err, "validate update data")
	}

	sportEvent, err := k.getSportEventToResolve(ctx, *resolutionPayload)
	if err != nil {
		return sdkerrors.Wrap(err, "getting sport-event")
	}

	if err := extractWinnerOddsIDs(&sportEvent, resolutionPayload); err != nil {
		return sdkerrors.Wrap(err, "extract winner odds id")
	}

	if err := k.Keeper.ResolveSportEvent(ctx, resolutionPayload); err != nil {
		return sdkerrors.Wrap(err, "resolve sport-event")
	}

	return nil
}

func (k msgServer) getSportEventToResolve(ctx sdk.Context, resolutionPayload types.SportEventResolutionTicketPayload) (types.SportEvent, error) {
	sportEvent, found := k.Keeper.GetSportEvent(ctx, resolutionPayload.UID)
	if !found {
		return types.SportEvent{}, types.ErrEventNotFound
	}

	if sportEvent.Status != types.SportEventStatus_SPORT_EVENT_STATUS_PENDING {
		return types.SportEvent{}, types.ErrEventIsNotPending
	}

	return sportEvent, nil
}

func extractWinnerOddsIDs(sportEvent *types.SportEvent, event *types.SportEventResolutionTicketPayload) error {
	if event.Status == types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED {
		if event.ResolutionTS < sportEvent.StartTS {
			return types.ErrResolutionTimeLessTnStart
		}

		validWinnerOdds := true
		for _, wid := range event.WinnerOddsUIDs {
			validWinnerOdds = false
			for _, o := range sportEvent.Odds {
				if o.UID == wid {
					validWinnerOdds = true
				}
			}
			if !validWinnerOdds {
				break
			}
		}

		if !validWinnerOdds {
			return types.ErrInvalidWinnerOdds
		}
	}

	return nil
}
