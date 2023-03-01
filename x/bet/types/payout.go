package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// CalculatePayoutProfit calculates the amount of payout profit portion according to bet odds value and amount
func CalculatePayoutProfit(oddsType OddsType, oddsVal string, amount sdk.Int) (sdk.Dec, error) {
	payout, err := calculatePayout(oddsType, oddsVal, amount)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// bettor profit is the subtracted amount of payout from bet amount
	profit := payout.Sub(amount.ToDec())

	return profit, nil
}

// calculatePayout calculates the amount of payout according to bet odds value and amount
func calculatePayout(oddsType OddsType, oddsVal string, amount sdk.Int) (sdk.Dec, error) {
	var oType OddsTypeI

	// assign corresponding type to the interface instance
	switch oddsType {
	case OddsType_ODDS_TYPE_DECIMAL:
		oType = new(decimalOdds)

	case OddsType_ODDS_TYPE_FRACTIONAL:
		oType = new(fractionalOdds)

	case OddsType_ODDS_TYPE_MONEYLINE:
		oType = new(moneylineOdds)

	default:
		return sdk.ZeroDec(), ErrInvalidOddsType
	}

	// total payout should be paid to bettor
	payout, err := oType.CalculatePayout(oddsVal, amount)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return payout, nil
}

// CalculateBetAmount calculates the amount of bet according to bet odds value and payout profit
func CalculateBetAmount(oddsType OddsType, oddsVal string, payoutProfit sdk.Dec) (sdk.Dec, error) {
	betAmount, err := calculateBetAmount(oddsType, oddsVal, payoutProfit)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return betAmount, nil
}

// calculateBetAmount calculates the amount of bet according to bet odds value and payoutProfit
func calculateBetAmount(oddsType OddsType, oddsVal string, payoutProfit sdk.Dec) (sdk.Dec, error) {
	var oType OddsTypeI

	// assign corresponding type to the interface instance
	switch oddsType {
	case OddsType_ODDS_TYPE_DECIMAL:
		oType = new(decimalOdds)

	case OddsType_ODDS_TYPE_FRACTIONAL:
		oType = new(fractionalOdds)

	case OddsType_ODDS_TYPE_MONEYLINE:
		oType = new(moneylineOdds)

	default:
		return sdk.ZeroDec(), ErrInvalidOddsType
	}

	// total payout should be paid to bettor
	betAmount, err := oType.CalculateBetAmount(oddsVal, payoutProfit)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return betAmount, nil
}
