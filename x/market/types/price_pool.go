package types

import sdkmath "cosmossdk.io/math"

func (pp PricePool) Remaining() sdkmath.Int {
	return pp.ResolutionFunds.Sub(pp.SpentFunds)
}
