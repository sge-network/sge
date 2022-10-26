package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/sportevent/types"
)

// SetSportEvent sets a specific sport event in the store
func (k Keeper) SetSportEvent(ctx sdk.Context, sportEvent types.SportEvent) {
	store := k.getSportEventsStore(ctx)
	b := k.cdc.MustMarshal(&sportEvent)
	store.Set(utils.StrBytes(sportEvent.UID), b)
}

// GetSportEvent returns a specific sport event by its UID
func (k Keeper) GetSportEvent(ctx sdk.Context, sportEventUID string) (val types.SportEvent, found bool) {
	sportEventsStore := k.getSportEventsStore(ctx)
	b := sportEventsStore.Get(utils.StrBytes(sportEventUID))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SportEventExists checks if a specific sport event exists or not
func (k Keeper) SportEventExists(ctx sdk.Context, sportEventUID string) bool {
	_, found := k.GetSportEvent(ctx, sportEventUID)
	return found
}

// RemoveSportEvent removes a sport event from the store
func (k Keeper) RemoveSportEvent(ctx sdk.Context, sportEventUID string) {
	store := k.getSportEventsStore(ctx)
	store.Delete(utils.StrBytes(sportEventUID))
}

// GetSportEventAll returns all sport events
func (k Keeper) GetSportEventAll(ctx sdk.Context) (list []types.SportEvent, err error) {
	store := k.getSportEventsStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SportEvent
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ResolveSportEvents updates a sport event with its resolution
func (k Keeper) ResolveSportEvents(ctx sdk.Context, resolutionEvents []types.ResolutionEvent) error {
	for _, resolutionEvent := range resolutionEvents {
		storedEvent, found := k.GetSportEvent(ctx, resolutionEvent.UID)
		if !found {
			return types.ErrNoMatchingSportEvent
		}

		if storedEvent.Status != types.SportEventStatus_STATUS_PENDING {
			return types.ErrCanNotBeAltered
		}

		storedEvent.Active = false
		storedEvent.ResolutionTS = resolutionEvent.ResolutionTS
		storedEvent.WinnerOddsUIDs = resolutionEvent.WinnerOddsUIDs
		storedEvent.Status = resolutionEvent.Status

		k.SetSportEvent(ctx, storedEvent)
	}
	return nil
}

// AddExtraPayoutToEvent update current total amount of payouts for a sport event
func (k Keeper) AddExtraPayoutToEvent(ctx sdk.Context, sportEventUID string, amount sdk.Int) error {
	sportEvent, found := k.GetSportEvent(ctx, sportEventUID)
	if !found {
		return types.ErrNoMatchingSportEvent
	}

	// calculate new total payout
	newAmount := sportEvent.BetConstraints.CurrentTotalAmount.Add(amount)
	if newAmount.GT(sportEvent.BetConstraints.MaxBetCap) {
		return types.ErrMaxBetCapExceeded
	}

	// update bet constraints of sport event in module state
	sportEvent.BetConstraints.CurrentTotalAmount = newAmount
	k.SetSportEvent(ctx, sportEvent)

	return nil
}

func emitTransactionEvent(ctx sdk.Context, emitType string, response *types.MsgSportResponse, creator string) {
	var failedString string
	if len(response.FailedEvents) != 0 {
		for i := range response.FailedEvents {
			failedString = failedString + response.FailedEvents[i].ID
			if i != len(response.FailedEvents)-1 {
				failedString = failedString + ", "
			}

		}
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			emitType,
			sdk.NewAttribute(types.AttributeKeySportEventsSuccessUID, strings.Join(response.SuccessEvents, ",")),
			sdk.NewAttribute(types.AttributeKeySportEventsFailedUID, failedString),
			sdk.NewAttribute(types.AttributeKeyEventsCreator, creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, emitType),
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
		),
	})
}
