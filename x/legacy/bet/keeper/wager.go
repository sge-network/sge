package keeper

import (
	"github.com/spf13/cast"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/legacy/bet/types"
	markettypes "github.com/sge-network/sge/x/legacy/market/types"
)

// Wager stores a new bet in KVStore
func (k Keeper) Wager(ctx sdk.Context, bet *types.Bet, betOdds map[string]*types.BetOddsCompact) error {
	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	market, err := k.getMarket(ctx, bet.MarketUID)
	if err != nil {
		return err
	}

	// check if selected odds is valid
	if !market.HasOdds(bet.OddsUID) {
		return types.ErrOddsUIDNotExist
	}

	if len(market.Odds) != len(betOdds) {
		return types.ErrInsufficientOdds
	}

	for _, odd := range market.Odds {
		if _, ok := betOdds[odd.UID]; !ok {
			return sdkerrors.Wrapf(types.ErrOddsUIDNotExist, "%s", odd.UID)
		}
	}

	// check minimum bet amount allowed
	betConstraints := k.GetConstraints(ctx)

	if bet.Amount.LT(betConstraints.MinAmount) {
		return types.ErrBetAmountIsLow
	}

	// modify the bet fee and subtracted amount
	bet.SetFee(betConstraints.Fee)

	// calculate payoutProfit
	payoutProfit, err := types.CalculatePayoutProfit(bet.OddsValue, bet.Amount)
	if err != nil {
		return err
	}

	stats := k.GetBetStats(ctx)
	stats.Count++
	betID := stats.Count

	betFulfillment, err := k.orderbookKeeper.ProcessWager(
		ctx, bet.UID, bet.MarketUID, bet.OddsUID, bet.MaxLossMultiplier, bet.Amount, payoutProfit,
		bettorAddress, bet.Fee, bet.OddsValue, betID, betOdds, market.OddsUIDS(),
	)
	if err != nil {
		return err
	}

	// set bet as placed
	bet.Status = types.Bet_STATUS_PLACED

	// put bet in the result pending status
	bet.Result = types.Bet_RESULT_PENDING

	bet.CreatedAt = ctx.BlockTime().Unix()
	bet.BetFulfillment = betFulfillment

	// store bet in the module state
	k.SetBet(ctx, *bet, betID)

	// set bet as a pending bet
	k.SetPendingBet(ctx, types.NewPendingBet(bet.UID, bet.Creator), betID, bet.MarketUID)

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
