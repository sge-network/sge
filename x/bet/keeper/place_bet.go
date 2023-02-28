package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
	"github.com/spf13/cast"
)

// PlaceBet stores a new bet in KVStore
func (k Keeper) PlaceBet(ctx sdk.Context, bet *types.Bet) error {
	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidCreator, err)
	}

	sportEvent, err := k.getSportEvent(ctx, bet.SportEventUID)
	if err != nil {
		return err
	}

	// check if selected odds is valid
	if !oddsExists(bet.OddsUID, sportEvent.Odds) {
		return types.ErrOddsUIDNotExist
	}

	// check minimum bet amount allowed
	betConstraints := sportEvent.GetBetConstraints()
	if betConstraints == nil {
		sportEvent.BetConstraints = k.sporteventKeeper.GetDefaultBetConstraints(ctx)
	}

	if bet.Amount.LT(sportEvent.BetConstraints.MinAmount) {
		return types.ErrBetAmountIsLow
	}

	// modify the bet fee and subtracted amount
	setBetFee(bet, sportEvent.BetConstraints.BetFee)

	// calculate payoutProfit
	payoutProfit, err := types.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
	if err != nil {
		return err
	}

	stats := k.GetBetStats(ctx)
	stats.Count++
	betID := stats.Count

	betFulfillment, err := k.obKeeper.ProcessBetPlacement(
		ctx, bet.UID, bet.SportEventUID, bet.OddsUID, bet.MaxLossMultiplier, payoutProfit,
		bettorAddress, bet.BetFee, bet.OddsType, bet.OddsValue, betID,
	)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInSRPlacementProcessing, "%s", err)
	}

	// set bet as placed
	bet.Status = types.Bet_STATUS_PLACED

	// put bet in the result pending status
	bet.Result = types.Bet_RESULT_PENDING

	bet.CreatedAt = ctx.BlockTime().Unix()
	bet.BetFulfillment = betFulfillment

	// store bet in the module state
	k.SetBet(ctx, *bet, betID)

	// set bet as an active bet
	k.SetActiveBet(ctx, types.NewActiveBet(bet.UID, bet.Creator), betID, bet.SportEventUID)

	// set bet stats
	k.SetBetStats(ctx, stats)

	return nil
}

// getSportEvent returns sport-event with id
func (k Keeper) getSportEvent(ctx sdk.Context, sportEventID string) (sporteventtypes.SportEvent, error) {
	sportevent, found := k.sporteventKeeper.GetSportEvent(ctx, sportEventID)
	if !found {
		return sporteventtypes.SportEvent{}, types.ErrNoMatchingSportEvent
	}

	if sportevent.Status != sporteventtypes.SportEventStatus_SPORT_EVENT_STATUS_ACTIVE {
		return sporteventtypes.SportEvent{}, types.ErrInactiveSportEvent
	}

	if sportevent.EndTS < cast.ToUint64(ctx.BlockTime().Unix()) {
		return sporteventtypes.SportEvent{}, types.ErrEndTSIsPassed
	}

	return sportevent, nil
}

// oddsExists checks if bet odds id is present in the sport-event list of odds uids
func oddsExists(betOddsUID string, odds []*sporteventtypes.Odds) bool {
	for _, o := range odds {
		if betOddsUID == o.UID {
			return true
		}
	}
	return false
}

// setBetFee sets the bet fee and subtraceted amount of bet object pointer
func setBetFee(bet *types.Bet, betFee sdk.Int) {
	bet.Amount = bet.Amount.Sub(betFee)
	bet.BetFee = betFee
}
