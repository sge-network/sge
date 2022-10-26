package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
)

// calculatePayout calculates the amount of payout according to bet odds value and amount
func calculatePayout(bet *types.Bet) sdk.Int {
	return (bet.OddsValue.MulInt(bet.Amount)).TruncateInt().Sub(bet.Amount)
}
