package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// CalculatePayoutProfit calculates the amount of payout profit portion according to bet odds value and amount
func CalculatePayoutProfit(oddsType OddsType, oddsVal string, amount sdk.Int) (sdk.Int, error) {
	payout, err := calculatePayout(oddsType, oddsVal, amount)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// bettor profit is the subtracted amount of payout from bet amount
	profit := payout.Sub(amount)

	return profit, nil
}

// calculatePayout calculates the amount of payout according to bet odds value and amount
func calculatePayout(oddsType OddsType, oddsVal string, amount sdk.Int) (sdk.Int, error) {
	var oType OddsTypeI

	// assign corresponding type to the interface instance
	switch oddsType {

	case OddsType_ODD_TYPE_DECIMAL:
		oType = new(decimalOdds)

	case OddsType_ODD_TYPE_FRACTIONAL:
		oType = new(fractionalOdds)

	case OddsType_ODD_TYPE_MONEYLINE:
		oType = new(moneylineOdds)

	default:
		return sdk.ZeroInt(), ErrInvalidOddsType

	}

	// total payout should be paid to bettor
	payout, err := oType.CalculatePayout(oddsVal, amount)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	return payout, nil
}
