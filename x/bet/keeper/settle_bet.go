package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
)

// SettleBet settles a single bet and updates it in KVStore
func (k Keeper) SettleBet(ctx sdk.Context, betUID string) error {

	if !types.IsValidUID(betUID) {
		return types.ErrInvalidBetUID
	}

	bet, found := k.GetBet(ctx, betUID)
	if !found {
		return types.ErrNoMatchingBet
	}

	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidCreator, err)
	}

	if err := checkBetStatus(bet.Status); err != nil {
		// bet cancelation logic will reside here if this feature is requested
		return err
	}

	// get the respective sport event for the bet
	sportEvent, found := k.sporteventKeeper.GetSportEvent(ctx, bet.SportEventUID)
	if !found {
		return types.ErrNoMatchingSportEvent
	}

	if sportEvent.Status == sporteventtypes.SportEventStatus_STATUS_ABORTED ||
		sportEvent.Status == sporteventtypes.SportEventStatus_STATUS_CANCELLED {
		bet.Result = types.Bet_RESULT_ABORTED

		payout := calculatePayout(&bet)
		if err := k.strategicreserveKeeper.RefundBettor(ctx, bettorAddress, bet.Amount, payout, bet.UID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRRefund, "%s", err)
		}

		bet.Status = types.Bet_STATUS_SETTLED

		k.SetBet(ctx, bet)
		return nil
	}

	// check if the bet odds is a winner odds or not and set the bet pointer states
	if err := resolveBetResult(&bet, sportEvent); err != nil {
		return err
	}

	if err := k.settle(ctx, &bet); err != nil {
		return err
	}

	// store bet in the module state
	k.SetBet(ctx, bet)
	return nil
}

// checkBetStatus checks status of bet. It returns an error if
// bet is canceled or settled alredy
func checkBetStatus(betstatus types.Bet_Status) error {

	switch betstatus {
	case types.Bet_STATUS_SETTLED:
		return types.ErrBetIsSettled
	case types.Bet_STATUS_CANCELLED:
		return types.ErrBetIsCanceled
	}

	return nil
}

// ResolveBetResult determines the result of the given bet, it can be lost or won.
func resolveBetResult(bet *types.Bet, sportEvent sporteventtypes.SportEvent) error {

	// check if sport event result is declared or not
	if sportEvent.Status != sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED {
		return types.ErrResultNotDeclared
	}

	_, exist := sportEvent.WinnerOddsUIDs[bet.OddsUID]
	if exist {
		// bet is winner
		bet.Result = types.Bet_RESULT_WON
		bet.Status = types.Bet_STATUS_RESULT_DECLARED
		return nil
	}

	// bet is loser
	bet.Result = types.Bet_RESULT_LOST
	bet.Status = types.Bet_STATUS_RESULT_DECLARED
	return nil
}

// settle settles a bet by calling strategicReserve functions to unlock fund and payout
// based on bet's result, and updates status of bet to settled
func (k Keeper) settle(ctx sdk.Context, bet *types.Bet) error {

	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidCreator, err)
	}

	payout := calculatePayout(bet)

	if bet.Result == types.Bet_RESULT_LOST {
		if err := k.strategicreserveKeeper.BettorLoses(ctx, bettorAddress, bet.Amount,
			payout, bet.UID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRBettorLoses, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED

	} else if bet.Result == types.Bet_RESULT_WON {
		if err := k.strategicreserveKeeper.BettorWins(ctx, bettorAddress, bet.Amount, payout, bet.UID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRBettorWins, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED
	}
	return nil
}
