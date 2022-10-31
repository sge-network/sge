package keeper

import (
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

	sportevent, err := k.getSportEvent(ctx, bet.SportEventUID)
	if err != nil {
		return err
	}

	if !oddsExists(bet.OddsUID, sportevent.OddsUIDs) {
		return types.ErrOddsUIDNotExist
	}

	// check minimum bet amount allowed
	if bet.Amount.LT(sportevent.BetConstraints.MinAmount) {
		return types.ErrBetAmountIsLow
	}

	vig := calculateVig(activeBetOdds)

	if vig.IsNegative() || vig.GT(sportevent.BetConstraints.MaxVig) || vig.LT(sportevent.BetConstraints.MinVig) {
		return types.ErrVigIsOutOfRange
	}

	// modify the bet fee and subtracted amount
	setBetFee(bet, sportevent.BetConstraints.BetFee)

	// calculate extraPayout
	extraPayout := calculateExtraPayout(bet)

	if err := k.sporteventKeeper.AddExtraPayoutToEvent(ctx, bet.SportEventUID, bet.OddsUID, bet.Amount, extraPayout); err != nil {
		return sdkerrors.Wrapf(types.ErrInAddAmountToSportEvent, "%s", err)
	}

	err = k.strategicreserveKeeper.ProcessBetPlacement(ctx, bettorAddress,
		bet.BetFee, bet.Amount, extraPayout, bet.UID)
	if err != nil {
		// bet placement was not successful so we need to update the total bet amount and payout statistics in the
		// sport event bet constraints
		if err := k.sporteventKeeper.AddExtraPayoutToEvent(ctx, bet.SportEventUID, bet.OddsUID, bet.Amount.Neg(), extraPayout.Neg()); err != nil {
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

	// emit events
	emitPlacementEvent(ctx, bet)

	return nil
}

func emitPlacementEvent(ctx sdk.Context, bet *types.Bet) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgPlaceBet,
			sdk.NewAttribute(types.AttributeKeyBetCreator, bet.Creator),
			sdk.NewAttribute(types.AttributeKeyBetUID, bet.UID),
			sdk.NewAttribute(types.AttributeKeySportEventUID, bet.SportEventUID),
			sdk.NewAttribute(sdk.AttributeKeyAmount, bet.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgPlaceBet),
			sdk.NewAttribute(sdk.AttributeKeySender, bet.Creator),
		),
	})
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

	if sportevent.EndTS < uint64(ctx.BlockTime().Unix()) {
		return sporteventtypes.SportEvent{}, types.ErrEndTSIsPassed
	}
	return sportevent, nil
}

// oddsExists checks if bet odds id is present in the sport event odds uids
func oddsExists(betOddsID string, sporteventOddsUIDs []string) bool {
	for _, oddsUID := range sporteventOddsUIDs {
		if betOddsID == oddsUID {
			return true
		}
	}
	return false
}

// calculates vig according to the active odds
func calculateVig(odds []*types.BetOdds) sdk.Dec {
	total := sdk.NewDec(0)
	for _, o := range odds {
		oVal := sdk.MustNewDecFromStr(o.Value)
		oValPerc := sdk.NewDec(100).Quo(oVal)
		total = total.Add(oValPerc)
	}

	return total.Sub(sdk.NewDec(types.MaxOddsValueSum))
}

// setBetFee sets the bet fee and subtraceted amount of bet object pointer
func setBetFee(bet *types.Bet, betFee sdk.Coin) {
	bet.Amount = bet.Amount.Sub(betFee.Amount)
	bet.BetFee = betFee
}
