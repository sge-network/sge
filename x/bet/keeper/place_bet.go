package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
)

// PlaceBet stores a new bet in KVStore
func (k Keeper) PlaceBet(ctx sdk.Context, bet *types.Bet, activeBetOdds []*types.BetOdds) error {
	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidCreator, err)
	}

	sportEvent, err := k.getSportEvent(ctx, bet.SportEventUID)
	if err != nil {
		return err
	}

	// check if selected odds is valid
	if !selectedOddsExists(bet.OddsUID, sportEvent.OddsUIDs) {
		return types.ErrOddsUIDNotExist
	}

	// check if provided active odds are valid
	if !allActiveOddsExist(activeBetOdds, bet.SportEventUID, sportEvent.OddsUIDs) {
		return types.ErrActiveOddsUIDsNotValid
	}

	// check minimum bet amount allowed
	if bet.Amount.LT(sportEvent.BetConstraints.MinAmount) {
		return types.ErrBetAmountIsLow
	}

	// modify the bet fee and subtracted amount
	setBetFee(bet, sportEvent.BetConstraints.BetFee)

	// calculate extraPayout
	extraPayout := calculateExtraPayout(bet)

	if err := k.sporteventKeeper.AddExtraPayoutToEvent(ctx, sportEvent.UID, extraPayout); err != nil {
		return sdkerrors.Wrapf(types.ErrInAddAmountToSportEvent, "%s", err)
	}

	err = k.strategicreserveKeeper.ProcessBetPlacement(ctx, bettorAddress,
		bet.BetFee, bet.Amount, extraPayout, bet.UID)
	if err != nil {
		// bet placement was not successful so we need to update the total bet amount and payout statistics in the
		// sport event bet constraints
		if err := k.sporteventKeeper.AddExtraPayoutToEvent(ctx, sportEvent.UID, extraPayout.Neg()); err != nil {
			return sdkerrors.Wrapf(types.ErrInSubAmountFromSportEvent, "%s", err)
		}
		return sdkerrors.Wrapf(types.ErrInSRPlacementProcessing, "%s", err)
	}

	// set bet as placed
	bet.Status = types.Bet_STATUS_PLACED

	// put bet in the result pending status
	bet.Result = types.Bet_RESULT_PENDING

	// set data verification status as valid
	bet.Verified = true

	bet.CreatedAt = ctx.BlockTime().Unix()

	// store bet in the module state
	k.SetBet(ctx, *bet)
	return nil
}

// getSportEvent returns sport event with id
func (k Keeper) getSportEvent(ctx sdk.Context, sportEventID string) (sporteventtypes.SportEvent, error) {
	sportevent, found := k.sporteventKeeper.GetSportEvent(ctx, sportEventID)
	if !found {
		return sporteventtypes.SportEvent{}, types.ErrNoMatchingSportEvent
	}

	if !sportevent.Active {
		return sporteventtypes.SportEvent{}, types.ErrInactiveSportEvent
	}

	if sportevent.Status != sporteventtypes.SportEventStatus_STATUS_PENDING {
		return sporteventtypes.SportEvent{}, types.ErrSportEventStatusNotPending
	}

	if sportevent.EndTS < uint64(time.Now().Unix()) {
		return sporteventtypes.SportEvent{}, types.ErrEndTSIsPassed
	}
	return sportevent, nil
}

// selectedOddsExists checks if bet odds id is present in the sport event odds uids
func selectedOddsExists(betOddsID string, sporteventOddsUIDs []string) bool {
	for _, oddsUID := range sporteventOddsUIDs {
		if betOddsID == oddsUID {
			return true
		}
	}
	return false
}

// allActiveOddsExist checks if all provided odds UIDs in activeBetOddss are related to the SportEventUID and are present in the sporteventOddsUIDs
func allActiveOddsExist(activeBetOdds []*types.BetOdds, SportEventUID string, sporteventOddsUIDs []string) bool {
	for _, odds := range activeBetOdds {
		if odds.SportEventUID != SportEventUID {
			return false
		}
	}
outerLoop:
	for _, providedOdds := range activeBetOdds {
		for _, sporteventOddsUID := range sporteventOddsUIDs {
			if providedOdds.UID == sporteventOddsUID {
				continue outerLoop
			}
		}
		return false
	}
	return true
}

// setBetFee sets the bet fee and subtraceted amount of bet object pointer
func setBetFee(bet *types.Bet, betFee sdk.Coin) {
	bet.Amount = bet.Amount.Sub(betFee.Amount)
	bet.BetFee = betFee
}
