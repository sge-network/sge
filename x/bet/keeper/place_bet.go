package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/spf13/cast"
)

// PlaceBet stores a new bet in KVStore
func (k Keeper) PlaceBet(ctx sdk.Context, bet *types.Bet) error {
	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidCreator, err)
	}

	market, err := k.getMarket(ctx, bet.MarketUID)
	if err != nil {
		return err
	}

	// check if selected odds is valid
	if !oddsExists(bet.OddsUID, market.Odds) {
		return types.ErrOddsUIDNotExist
	}

	// check minimum bet amount allowed
	betConstraints := market.GetBetConstraints()
	if betConstraints == nil {
		market.BetConstraints = k.marketKeeper.GetDefaultBetConstraints(ctx)
	}

	if bet.Amount.LT(market.BetConstraints.MinAmount) {
		return types.ErrBetAmountIsLow
	}

	// modify the bet fee and subtracted amount
	setBetFee(bet, market.BetConstraints.BetFee)

	// calculate payoutProfit
	payoutProfit, err := types.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
	if err != nil {
		return err
	}

	stats := k.GetBetStats(ctx)
	stats.Count++
	betID := stats.Count

	betFulfillment, err := k.srKeeper.ProcessBetPlacement(
		ctx, bet.UID, bet.MarketUID, bet.OddsUID, bet.MaxLossMultiplier, bet.Amount, payoutProfit,
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
	k.SetActiveBet(ctx, types.NewActiveBet(bet.UID, bet.Creator), betID, bet.MarketUID)

	// set bet stats
	k.SetBetStats(ctx, stats)

	return nil
}

// getMarket returns market with id
func (k Keeper) getMarket(ctx sdk.Context, marketID string) (markettypes.Market, error) {
	market, found := k.marketKeeper.GetMarket(ctx, marketID)
	if !found {
		return markettypes.Market{}, types.ErrNoMatchingMarket
	}

	if market.Status != markettypes.MarketStatus_MARKET_STATUS_ACTIVE {
		return markettypes.Market{}, types.ErrInactiveMarket
	}

	if market.EndTS < cast.ToUint64(ctx.BlockTime().Unix()) {
		return markettypes.Market{}, types.ErrEndTSIsPassed
	}

	return market, nil
}

// oddsExists checks if bet odds id is present in the market list of odds uids
func oddsExists(betOddsUID string, odds []*markettypes.Odds) bool {
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
