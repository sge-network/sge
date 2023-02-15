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
	orderBookId, found := k.GetFirstUnsettledResolvedOrderBook(ctx)

	// return if there is no resolved orderbook.
	if !found {
		return nil
	}

	book, found := k.GetBook(ctx, orderBookId)
	if !found {
		return fmt.Errorf("orderbook not found %s", orderBookId)
	}
	if book.Status != types.OrderBookStatus_STATUS_RESOLVED {
		return fmt.Errorf("orderbook status not resolved %s", orderBookId)
	}

	// settle order book active deposits.
	allSettled, err := k.batchSettlementOfDeposit(ctx, orderBookId, toFetch)
	if err != nil {
		return fmt.Errorf("could not settle orderbook %s %s", orderBookId, err)
	}

	// if there is not any active deposit for orderbook
	// we need to remove its uid from the list of unsettled resolved orderbooks.
	if allSettled {
		k.RemoveUnsettledResolvedOrderBook(ctx, orderBookId)
	}

	return nil
}

// batchSettlementOfDeposit settles active deposits of a orderbook
func (k Keeper) batchSettlementOfDeposit(ctx sdk.Context, orderBookId string, countToBeSettled uint64) (allSettled bool, err error) {
	// initialize iterator for the certain number of active deposits
	// equal to countToBeSettled
	allSettled, settled := true, 0
	bookParticipants := k.GetParticipantsByBook(ctx, orderBookId)
	for _, bookParticipant := range bookParticipants {
		if !bookParticipant.IsSettled {
			err = k.settleDeposit(ctx, bookParticipant)
			if err != nil {
				return allSettled, err
			}
			settled += 1
			allSettled = false
		}
		if settled >= int(countToBeSettled) {
			break
		}
	}

	return allSettled, nil
}

func (k Keeper) settleDeposit(ctx sdk.Context, bp types.BookParticipant) error {
	if bp.IsSettled {
		return sdkerrors.Wrapf(types.ErrBookParticipantAlreadySettled, "%s %d", bp.BookId, bp.ParticipantNumber)
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
	k.SetBookParticipant(ctx, bp)
	return nil
}
