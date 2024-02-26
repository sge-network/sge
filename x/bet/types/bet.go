package types

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

func NewPendingBet(uid, creator string) *PendingBet {
	return &PendingBet{
		UID:     uid,
		Creator: creator,
	}
}

func NewSettledBet(uid, bettorAddress string) *SettledBet {
	return &SettledBet{
		UID:           uid,
		BettorAddress: bettorAddress,
	}
}

func NewBetFulfillment(
	participantAddress string,
	participationIndex uint64,
	betAmount, payoutProfit sdkmath.Int,
) *BetFulfillment {
	return &BetFulfillment{
		ParticipantAddress: participantAddress,
		ParticipationIndex: participationIndex,
		BetAmount:          betAmount,
		PayoutProfit:       payoutProfit,
	}
}

// CheckSettlementEligiblity checks status of bet. It returns an error if
// bet is canceled or settled already.
func (bet *Bet) CheckSettlementEligiblity() error {
	switch bet.Status {
	case Bet_STATUS_SETTLED:
		return ErrBetIsSettled
	case Bet_STATUS_CANCELED:
		return ErrBetIsCanceled
	}

	return nil
}

// SetResult sets the bet object results according to the market resolution.
func (bet *Bet) SetResult(market *markettypes.Market) error {
	// check if market result is declared or not
	if market.Status != markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED {
		return ErrResultNotDeclared
	}

	var exist bool
	for _, wid := range market.WinnerOddsUIDs {
		if wid == bet.OddsUID {
			exist = true
			break
		}
	}

	if exist {
		// bettor is winner
		bet.Result = Bet_RESULT_WON
		bet.Status = Bet_STATUS_RESULT_DECLARED
		return nil
	}

	// bettor is loser
	bet.Result = Bet_RESULT_LOST
	bet.Status = Bet_STATUS_RESULT_DECLARED
	return nil
}

// SetFee calculates and sets the betting fee.
func (bet *Bet) SetFee(fee sdkmath.Int) {
	bet.Amount = bet.Amount.Sub(fee)
	bet.Fee = fee
}

// SetPriceLockFee calculates and sets the price lock fee.
func (bet *Bet) SetPriceLockFee(feePercentage sdk.Dec) {
	feeAmount := sdk.NewDecFromInt(bet.Amount).Mul(feePercentage).TruncateInt()
	bet.PriceLockFee = feeAmount
	bet.Amount = bet.Amount.Sub(feeAmount)
}

// IsPriceLockEnabled returns true if the price lock feature is requested for this bet
func (bet *Bet) IsPriceLockEnabled() bool {
	return !bet.WagerSgePrice.IsNil() && bet.WagerSgePrice.GT(sdk.ZeroDec())
}

// SetPriceReimbursement calculates and sets the price reimbursement.
func (bet *Bet) SetPriceReimbursement(resolutionPrice sdk.Dec) {
	if bet.WagerSgePrice.IsZero() {
		return
	}

	for _, bf := range bet.BetFulfillment {
		if bet.PriceReimbursement.IsNil() {
			bet.PriceReimbursement = sdkmath.ZeroInt()
		}
		bet.PriceReimbursement = bet.PriceReimbursement.Add(
			bf.PriceReimbursed(
				bet.WagerSgePrice,
				resolutionPrice,
			),
		)
	}
}

// PriceReimbursed calculates the price reimbursement amount.
func (bf *BetFulfillment) PriceReimbursed(first, second sdk.Dec) sdkmath.Int {
	if second.LT(first) {
		// sge token value dropped
		amountAndProfit := bf.BetAmount.Add(bf.PayoutProfit)
		firstValueInDollars := sdk.NewDecFromInt(amountAndProfit).Mul(first)
		secondValueInDollars := sdk.NewDecFromInt(amountAndProfit).Mul(second)
		// sge reimbursement tokens
		reimbursementAmount := firstValueInDollars.Sub(secondValueInDollars)
		return reimbursementAmount.TruncateInt()
	}

	return sdkmath.ZeroInt()
}
