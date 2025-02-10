package keeper

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	storetypes "cosmossdk.io/store/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/legacy/bet/types"
	markettypes "github.com/sge-network/sge/x/legacy/market/types"
)

// singlePageNum used to return single page result in pagination.
const singlePageNum = 1

// Settle settles a single bet and updates it in KVStore
func (k Keeper) Settle(ctx sdk.Context, bettorAddressStr, betUID string) error {
	if !utils.IsValidUID(betUID) {
		return types.ErrInvalidBetUID
	}

	uid2ID, found := k.GetBetID(ctx, betUID)
	if !found {
		return types.ErrNoMatchingBet
	}

	bet, found := k.GetBet(ctx, bettorAddressStr, uid2ID.ID)
	if !found {
		return types.ErrNoMatchingBet
	}

	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	if bet.Creator != bettorAddressStr {
		return types.ErrBettorAddressNotEqualToCreator
	}

	if err := bet.CheckSettlementEligiblity(); err != nil {
		// bet cancellation logic will reside here if this feature is requested
		return err
	}

	// get the respective market for the bet
	market, found := k.marketKeeper.GetMarket(ctx, bet.MarketUID)
	if !found {
		return types.ErrNoMatchingMarket
	}

	if market.Status == markettypes.MarketStatus_MARKET_STATUS_ABORTED ||
		market.Status == markettypes.MarketStatus_MARKET_STATUS_CANCELED {
		payoutProfit, err := types.CalculatePayoutProfit(bet.OddsValue, bet.Amount)
		if err != nil {
			return err
		}

		if err := k.orderbookKeeper.RefundBettor(ctx, bettorAddress, bet.Amount, bet.Fee, payoutProfit.TruncateInt(), bet.UID); err != nil {
			return sdkerrors.Wrapf(types.ErrInOBRefund, "%s", err)
		}

		bet.Status = types.Bet_STATUS_SETTLED
		bet.Result = types.Bet_RESULT_REFUNDED

		k.updateSettlementState(ctx, bet, uid2ID.ID)

		return nil
	}

	// check if the bet odds is a winner odds or not and set the bet pointer states
	if err := bet.SetResult(&market); err != nil {
		return err
	}

	if err := k.settleResolved(ctx, &bet); err != nil {
		return err
	}

	if err := k.orderbookKeeper.WithdrawBetFee(ctx, sdk.MustAccAddressFromBech32(market.Creator), bet.Fee); err != nil {
		return err
	}

	k.updateSettlementState(ctx, bet, uid2ID.ID)

	return nil
}

// updateSettlementState settles bet in the store
func (k Keeper) updateSettlementState(ctx sdk.Context, bet types.Bet, betID uint64) {
	// set current height as settlement height
	bet.SettlementHeight = ctx.BlockHeight()

	// store bet in the module state
	k.SetBet(ctx, bet, betID)

	// remove pending bet
	k.RemovePendingBet(ctx, bet.MarketUID, betID)

	// store settled bet in the module state
	k.SetSettledBet(ctx, types.NewSettledBet(bet.UID, bet.Creator), betID, ctx.BlockHeight())
}

// settleResolved settles a bet by calling order book functions to unlock fund and payout
// based on bet's result, and updates status of bet to settled
func (k Keeper) settleResolved(ctx sdk.Context, bet *types.Bet) error {
	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	if bet.Result == types.Bet_RESULT_LOST {
		if err := k.orderbookKeeper.BettorLoses(ctx, bet.BetFulfillment, bet.MarketUID); err != nil {
			return sdkerrors.Wrapf(types.ErrInOBBettorLoses, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED
	} else if bet.Result == types.Bet_RESULT_WON {
		if err := k.orderbookKeeper.BettorWins(ctx, bettorAddress, bet.BetFulfillment, bet.MarketUID); err != nil {
			return sdkerrors.Wrapf(types.ErrInOBBettorWins, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED
	}
	return nil
}

// BatchMarketSettlements settles bets of resolved markets
// in batch. The markets get into account in FIFO manner.
func (k Keeper) BatchMarketSettlements(ctx sdk.Context) error {
	toFetch := k.GetParams(ctx).BatchSettlementCount

	// continue looping until reach batch settlement count parameter
	for toFetch > 0 {
		// get the first resolved market to process corresponding pending bets.
		marketUID, found := k.marketKeeper.GetFirstUnsettledResolvedMarket(ctx)
		// exit loop if there is no resolved bet.
		if !found {
			return nil
		}

		// settle market pending bets.
		settledCount, err := k.batchMarketSettlement(ctx, marketUID, toFetch)
		if err != nil {
			return fmt.Errorf("could not settle market %s %s", marketUID, err)
		}

		// check if still there is any pending bet for the market.
		pendingBetExists, err := k.IsAnyPendingBetForMarket(ctx, marketUID)
		if err != nil {
			return fmt.Errorf("could not check the pending bets %s %s", marketUID, err)
		}

		// if there is not any pending bet for the market
		// we need to remove its uid from the list of unsettled resolved bets.
		if !pendingBetExists {
			k.marketKeeper.RemoveUnsettledResolvedMarket(ctx, marketUID)
			err = k.orderbookKeeper.SetOrderBookAsUnsettledResolved(ctx, marketUID)
			if err != nil {
				return fmt.Errorf("could not resolve orderbook %s %s", marketUID, err)
			}
		}

		// update counter of bets to be processed in the next iteration.
		toFetch -= settledCount
	}

	return nil
}

// batchMarketSettlement settles pending bets of a markets
//
//nolint:nakedret
func (k Keeper) batchMarketSettlement(
	ctx sdk.Context,
	marketUID string,
	countToBeSettled uint32,
) (settledCount uint32, err error) {
	// initialize iterator for the certain number of pending bets
	// equal to countToBeSettled
	iterator := storetypes.KVStorePrefixIteratorPaginated(
		ctx.KVStore(k.storeKey),
		types.PendingBetListOfMarketPrefix(marketUID),
		singlePageNum,
		uint(countToBeSettled))
	defer func() {
		iterErr := iterator.Close()
		if iterErr != nil {
			err = iterErr
		}
	}()

	// settle bets for the filtered pending bets
	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingBet
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		err = k.Settle(ctx, val.Creator, val.UID)
		if err != nil {
			return
		}

		// update total settled count
		settledCount++
	}

	return settledCount, nil
}
