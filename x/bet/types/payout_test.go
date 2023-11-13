package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

var ecceptedTruncatedValue = sdk.NewInt(1)

var defaultBetAmount = int64(35625789)

func TestCalculateDecimalPayout(t *testing.T) {
	tcs := []struct {
		desc      string
		oddsValue string
		betAmount int64

		expVal int64
		err    error
	}{
		{
			desc:      "positive",
			oddsValue: "1.55",
			betAmount: defaultBetAmount,

			expVal: 19594183,
		},
		{
			desc:      "same",
			oddsValue: "1",
			betAmount: defaultBetAmount,

			err: types.ErrDecimalOddsCanNotBeLessThanOne,
		},
		{
			desc:      "negative",
			oddsValue: "0.75",
			betAmount: defaultBetAmount,

			err: types.ErrDecimalOddsCanNotBeLessThanOne,
		},
		{
			desc:      "negative input",
			oddsValue: "-0.75",
			betAmount: defaultBetAmount,

			err: types.ErrDecimalOddsShouldBePositive,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			payoutProfit, err := types.CalculatePayoutProfit(
				tc.oddsValue,
				sdk.NewInt(tc.betAmount),
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.expVal).Equal(payoutProfit.TruncateInt()), "expected: %d, actual: %d", tc.expVal, payoutProfit)
			}

			calcBetAmount, err := types.CalculateBetAmount(
				tc.oddsValue,
				payoutProfit,
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.betAmount).Sub(calcBetAmount.Ceil().TruncateInt()).LT(ecceptedTruncatedValue), "expected: %d, actual: %d", tc.betAmount, calcBetAmount)
			}
		})
	}
}
