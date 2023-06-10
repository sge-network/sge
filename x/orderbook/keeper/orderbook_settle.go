package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"

	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// BatchOrderBookSettlements settles order books
func (k Keeper) BatchOrderBookSettlements(ctx sdk.Context) error {
	toFetch := k.GetParams(ctx).BatchSettlementCount

	// get the first resolved orderboook to process corresponding active deposits.
	orderBookUID, found := k.GetFirstUnsettledResolvedOrderBook(ctx)

	// return if there is no resolved orderboook.
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
	allSettled, err := k.batchSettlementOfParticipation(ctx, orderBookUID, market, toFetch)
	if err != nil {
		return fmt.Errorf("could not settle orderboook %s %s", orderBookUID, err)
	}

	// if there is not any active deposit for orderboook
	// we need to remove its uid from the list of unsettled resolved orderbooks.
	if allSettled {
		k.RemoveUnsettledResolvedOrderBook(ctx, orderBookUID)
	}

	return nil
}

// batchSettlementOfParticipation settles active deposits of an orderboook
func (k Keeper) batchSettlementOfParticipation(
	ctx sdk.Context,
	orderBookUID string,
	market markettypes.Market,
	countToBeSettled uint64,
) (allSettled bool, err error) {
	// initialize iterator for the certain number of active deposits
	// equal to countToBeSettled
	allSettled, settled := true, 0
	bookParticipations, err := k.GetParticipationsOfOrderBook(ctx, orderBookUID)
	if err != nil {
		return false, fmt.Errorf("batch settlement of book %s failed: %s", orderBookUID, err)
	}
	for _, bookParticipation := range bookParticipations {
		if !bookParticipation.IsSettled {
			err = k.settleParticipation(ctx, bookParticipation, market)
			if err != nil {
				return allSettled, fmt.Errorf(
					"failed to settle deposit of batch settlement for participation %#v: %s",
					bookParticipation,
					err,
				)
			}
			settled++
			allSettled = false
		}
		if cast.ToUint64(settled) >= countToBeSettled {
			break
		}
	}

	return allSettled, nil
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidDepositor, err)
	}

	refundHouseDepositFeeToDepositor := false

	switch market.Status {
	case markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED:
		depositPlusProfit := bp.Liquidity.Add(bp.ActualProfit)
		// refund participant's account from orderbook liquidity pool.
		if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, depositorAddress, depositPlusProfit); err != nil {
			return err
		}
		if bp.NotParticipatedInBetFulfillment() {
			refundHouseDepositFeeToDepositor = true
		}
	case markettypes.MarketStatus_MARKET_STATUS_CANCELED,
		markettypes.MarketStatus_MARKET_STATUS_ABORTED:
		// refund participant's account from orderbook liquidity pool.
		if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, depositorAddress, bp.Liquidity); err != nil {
			return err
		}
		refundHouseDepositFeeToDepositor = true
	default:
		return sdkerrors.Wrapf(
			types.ErrUnknownMarketStatus,
			"order book %s,  market status %s",
			bp.OrderBookUID,
			market.Status,
		)
	}

	// get corresponding deposit to extract house fee
	deposit, found := k.houseKeeper.GetDeposit(ctx, bp.ParticipantAddress, bp.OrderBookUID, bp.Index)
	if !found {
		return sdkerrors.Wrapf(types.ErrDepositNotFoundForParticipation, "%s", err)
	}

	if refundHouseDepositFeeToDepositor {
		// refund participant's account from house fee collector.
		if err := k.refund(housetypes.HouseFeeCollectorFunder{}, ctx, depositorAddress, deposit.Fee); err != nil {
			return err
		}
	} else {
		// refund participant's account from house fee collector.
		if err := k.refund(housetypes.HouseFeeCollectorFunder{}, ctx, sdk.MustAccAddressFromBech32(market.Creator), deposit.Fee); err != nil {
			return err
		}
	}

	bp.IsSettled = true
	k.SetOrderBookParticipation(ctx, bp)
	return nil
}
