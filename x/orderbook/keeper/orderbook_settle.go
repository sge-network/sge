package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/orderbook/types"
	srtypes "github.com/sge-network/sge/x/strategicreserve/types"
)

// BatchOrderBookSettlements settles order books
func (k Keeper) BatchOrderBookSettlements(ctx sdk.Context) error {
	toFetch := k.GetParams(ctx).BatchSettlementCount

	// get the first resolved orderbook to process corresponding active deposits.
	orderBookUID, found := k.GetFirstUnsettledResolvedOrderBook(ctx)

	// return if there is no resolved orderbook.
	if !found {
		return nil
	}

	book, found := k.GetBook(ctx, orderBookUID)
	if !found {
		return fmt.Errorf("orderbook not found %s", orderBookUID)
	}
	if book.Status != types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_RESOLVED {
		return fmt.Errorf("orderbook status not resolved %s", orderBookUID)
	}

	// settle order book active deposits.
	allSettled, err := k.batchSettlementOfDeposit(ctx, orderBookUID, toFetch)
	if err != nil {
		return fmt.Errorf("could not settle orderbook %s %s", orderBookUID, err)
	}

	// if there is not any active deposit for orderbook
	// we need to remove its uid from the list of unsettled resolved orderbooks.
	if allSettled {
		k.RemoveUnsettledResolvedOrderBook(ctx, orderBookUID)
	}

	return nil
}

// batchSettlementOfDeposit settles active deposits of a orderbook
func (k Keeper) batchSettlementOfDeposit(ctx sdk.Context, orderBookUID string, countToBeSettled uint64) (allSettled bool, err error) {
	// initialize iterator for the certain number of active deposits
	// equal to countToBeSettled
	allSettled, settled := true, 0
	bookParticipations, err := k.GetParticipationsOfBook(ctx, orderBookUID)
	if err != nil {
		return false, fmt.Errorf("batch settlement of book %s failed: %s", orderBookUID, err)
	}
	for _, bookParticipation := range bookParticipations {
		if !bookParticipation.IsSettled {
			err = k.settleDeposit(ctx, bookParticipation)
			if err != nil {
				return allSettled, fmt.Errorf("failed to settle deposit of batch settlement for participation %#v: %s",
					bookParticipation, err)
			}
			settled++
			allSettled = false
		}
		if settled >= int(countToBeSettled) {
			break
		}
	}

	return allSettled, nil
}

func (k Keeper) settleDeposit(ctx sdk.Context, bp types.BookParticipation) error {
	if bp.IsSettled {
		return sdkerrors.Wrapf(types.ErrBookParticipationAlreadySettled, "%s %d", bp.BookUID, bp.Index)
	}

	if bp.IsModuleAccount {
		depositPlusProfit := bp.Liquidity.Add(bp.ActualProfit)
		if depositPlusProfit.LTE(bp.Liquidity) {
			// transfer amount to `sr_pool` module account
			err := k.transferFundsFromModuleToModule(ctx, types.BookLiquidityName, srtypes.SRPoolName, depositPlusProfit)
			if err != nil {
				return err
			}
		} else {
			// transfer initial amount to `sr_pool` module account
			err := k.transferFundsFromModuleToModule(ctx, types.BookLiquidityName, srtypes.SRPoolName, bp.Liquidity)
			if err != nil {
				return err
			}

			// transfer profit to `sr_profit_pool` module account
			err = k.transferFundsFromModuleToModule(ctx, types.BookLiquidityName, types.SRProfitName, bp.ActualProfit)
			if err != nil {
				return err
			}
		}
	} else {
		depositPlusProfit := bp.Liquidity.Add(bp.ActualProfit)
		depositorAddress, err := sdk.AccAddressFromBech32(bp.ParticipantAddress)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidDesositor, err)
		}
		if depositPlusProfit.LTE(bp.Liquidity) {
			// transfer amount to depositor address
			err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, depositorAddress, depositPlusProfit)
			if err != nil {
				return err
			}
		} else {
			// transfer initial amount to depositor address
			err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, depositorAddress, bp.Liquidity)
			if err != nil {
				return err
			}

			// transfer profit to depositor address
			err = k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, depositorAddress, bp.ActualProfit)
			if err != nil {
				return err
			}
		}
	}
	bp.IsSettled = true
	k.SetBookParticipation(ctx, bp)
	return nil
}
