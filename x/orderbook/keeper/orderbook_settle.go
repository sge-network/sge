package keeper

import (
	"fmt"

	"github.com/spf13/cast"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// BatchOrderBookSettlements settles order books
func (k Keeper) BatchOrderBookSettlements(ctx sdk.Context) error {
	toFetch := k.GetParams(ctx).BatchSettlementCount

	unresolvedOrderBookIndex := 0
	for toFetch > 0 {
		// get the first resolved orderbook to process corresponding active deposits.
		orderBookUID, found := k.GetFirstUnsettledResolvedOrderBook(ctx, unresolvedOrderBookIndex)

		// return if there is no resolved orderbook.
		if !found {
			return nil
		}

		book, found := k.GetOrderBook(ctx, orderBookUID)
		if !found {
			return fmt.Errorf("orderbook not found %s", orderBookUID)
		}
		if book.Status != types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_RESOLVED {
			return fmt.Errorf("orderbook status not resolved %s", orderBookUID)
		}

		market, found := k.marketKeeper.GetMarket(ctx, orderBookUID)
		if !found {
			return fmt.Errorf("market not found %s", orderBookUID)
		}

		// settle order book active deposits.
		allSettled, settledCount, err := k.batchSettlementOfParticipation(ctx, orderBookUID, market, toFetch)
		if err != nil {
			return fmt.Errorf("could not settle orderbook %s %s", orderBookUID, err)
		}

		// if there is not any active deposit for orderbook
		// we need to remove its uid from the list of unsettled resolved orderbooks.
		if allSettled {
			k.RemoveUnsettledResolvedOrderBook(ctx, orderBookUID)

			book.Status = types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_SETTLED
			k.SetOrderBook(ctx, book)
		}

		// update counter of bets to be processed in the next iteration.
		toFetch -= settledCount
		// update market index to be checked in the next loop.
		unresolvedOrderBookIndex++
	}

	return nil
}

// batchSettlementOfParticipation settles active deposits of an orderbook
func (k Keeper) batchSettlementOfParticipation(
	ctx sdk.Context,
	orderBookUID string,
	market markettypes.Market,
	countToBeSettled uint64,
) (allSettled bool, settledCount uint64, err error) {
	// initialize iterator for the certain number of active deposits
	// equal to countToBeSettled
	allSettled = true
	bookParticipations, err := k.GetParticipationsOfOrderBook(ctx, orderBookUID)
	if err != nil {
		return false, settledCount, fmt.Errorf("batch settlement of book %s failed: %s", orderBookUID, err)
	}
	for _, bookParticipation := range bookParticipations {
		if !bookParticipation.IsSettled {
			err = k.settleParticipation(ctx, bookParticipation, market)
			if err != nil {
				return allSettled, settledCount, fmt.Errorf(
					"failed to settle deposit of batch settlement for participation %#v: %s",
					bookParticipation,
					err,
				)
			}
			settledCount++
			allSettled = false
		}
		if cast.ToUint64(settledCount) >= countToBeSettled {
			break
		}
	}

	return allSettled, settledCount, nil
}

func (k Keeper) settleParticipation(
	ctx sdk.Context,
	bp types.OrderBookParticipation,
	market markettypes.Market,
) error {
	if bp.IsSettled {
		return sdkerrors.Wrapf(
			types.ErrBookParticipationAlreadySettled,
			"%s %d",
			bp.OrderBookUID,
			bp.Index,
		)
	}

	depositorAddress, err := sdk.AccAddressFromBech32(bp.ParticipantAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, types.ErrTextInvalidDepositor, err)
	}

	refundHouseDepositFeeToDepositor := false

	switch market.Status {
	case markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED:
		bp.ReturnedAmount = bp.Liquidity.Add(bp.ActualProfit)
		// refund participant's account from orderbook liquidity pool.
		if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, depositorAddress, bp.ReturnedAmount); err != nil {
			return err
		}
		if bp.NotParticipatedInBetFulfillment() {
			refundHouseDepositFeeToDepositor = true
		}
		if bp.ActualProfit.IsNegative() {
			k.hooks.AfterHouseLoss(ctx, depositorAddress, bp.Liquidity, bp.ActualProfit.Abs())
		} else {
			k.hooks.AfterHouseWin(ctx, depositorAddress, bp.Liquidity, bp.ActualProfit)
		}

	case markettypes.MarketStatus_MARKET_STATUS_CANCELED,
		markettypes.MarketStatus_MARKET_STATUS_ABORTED:
		bp.ReturnedAmount = bp.Liquidity
		// refund participant's account from orderbook liquidity pool.
		if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, depositorAddress, bp.ReturnedAmount); err != nil {
			return err
		}
		refundHouseDepositFeeToDepositor = true
		k.hooks.AfterHouseRefund(ctx, depositorAddress, bp.ReturnedAmount)
	default:
		return sdkerrors.Wrapf(
			types.ErrUnknownMarketStatus,
			"order book %s,  market status %s",
			bp.OrderBookUID,
			market.Status,
		)
	}

	if refundHouseDepositFeeToDepositor {
		// refund participant's account from house fee collector.
		if err := k.refund(housetypes.HouseFeeCollectorFunder{}, ctx, depositorAddress, bp.Fee); err != nil {
			return err
		}
		bp.ReimbursedFee = bp.Fee
		bp.ReturnedAmount = bp.ReturnedAmount.Add(bp.ReimbursedFee)
		k.hooks.AfterHouseFeeRefund(ctx, depositorAddress, bp.Fee)
	} else {
		// refund participant's account from house fee collector.
		if err := k.refund(housetypes.HouseFeeCollectorFunder{}, ctx, sdk.MustAccAddressFromBech32(market.Creator), bp.Fee); err != nil {
			return err
		}
	}

	// market uid
	// participation index
	// returned fees
	// returned liquidity
	// actual profit

	bp.IsSettled = true
	k.SetOrderBookParticipation(ctx, bp)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeParticipationSettlement,
			sdk.NewAttribute(types.AttributeValueMarketUID, market.UID),
			sdk.NewAttribute(types.AttributeValueParticipationIndex, cast.ToString(bp.Index)),
			sdk.NewAttribute(types.AttributeValueParticipationReimbursedFees, bp.ReimbursedFee.String()),
			sdk.NewAttribute(types.AttributeValueParticipationReturnedAmount, bp.ReturnedAmount.String()),
			sdk.NewAttribute(types.AttributeValueParticipationActualProfit, bp.ActualProfit.String()),
		),
	)

	return nil
}
