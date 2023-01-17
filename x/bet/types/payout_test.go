package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

var defaultBetAmount = int64(100)

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

			expVal: 55,
		},
		{
			desc:      "negative",
			oddsValue: "0.75",
			betAmount: defaultBetAmount,

			expVal: -25,
		},
		{
			desc:      "same",
			oddsValue: "1",
			betAmount: defaultBetAmount,

			expVal: 0,
		},
		{
			desc:      "negative input",
			oddsValue: "-0.75",
			betAmount: defaultBetAmount,

			err: types.ErrDecimalOddsCanNotBeNegative,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			payout, err := types.CalculatePayoutProfit(types.OddsType_ODD_TYPE_DECIMAL, tc.oddsValue, sdk.NewInt(tc.betAmount))
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.expVal).Equal(payout), "expected: %d, actual: %d", tc.expVal, payout.Int64())
			}
		})
	}
}

func TestCalculateFractionalPayout(t *testing.T) {

	tcs := []struct {
		desc      string
		oddsValue string
		betAmount int64

		expVal int64
		err    error
	}{
		{
			desc:      "positive outcome",
			oddsValue: "5/2",
			betAmount: defaultBetAmount,

			expVal: 250,
		},
		{
			desc:      "negative outcome",
			oddsValue: "2/7",
			betAmount: defaultBetAmount,

			expVal: 29,
		},
		{
			desc:      "same",
			oddsValue: "1/1",
			betAmount: defaultBetAmount,

			expVal: defaultBetAmount,
		},
		{
			desc:      "zero base",
			oddsValue: "1/0",
			betAmount: defaultBetAmount,

			err: types.ErrFractionalOddsCanNotBeNegativeOrZero,
		},
		{
			desc:      "negative base",
			oddsValue: "-5/2",
			betAmount: defaultBetAmount,

			err: types.ErrFractionalOddsCanNotBeNegativeOrZero,
		},
		{
			desc:      "zero top",
			oddsValue: "0/1",
			betAmount: defaultBetAmount,

			err: types.ErrFractionalOddsCanNotBeNegativeOrZero,
		},
		{
			desc:      "negative top",
			oddsValue: "5/-2",
			betAmount: defaultBetAmount,

			err: types.ErrFractionalOddsCanNotBeNegativeOrZero,
		},
		{
			desc:      "negative top",
			oddsValue: "5//2",
			betAmount: defaultBetAmount,

			err: types.ErrFractionalOddsIncorrectFormat,
		},
		{
			desc:      "wrong type",
			oddsValue: "5/s",
			betAmount: defaultBetAmount,

			err: types.ErrInConvertingOddsToInt,
		},
		{
			desc:      "incorrect format",
			oddsValue: "5//2",
			betAmount: defaultBetAmount,

			err: types.ErrFractionalOddsIncorrectFormat,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			payout, err := types.CalculatePayoutProfit(types.OddsType_ODD_TYPE_FRACTIONAL, tc.oddsValue, sdk.NewInt(tc.betAmount))
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.expVal).Equal(payout), "expected: %d, actual: %d", tc.expVal, payout.Int64())
			}
		})
	}

}

func TestCalculateMoneylinePayout(t *testing.T) {

	tcs := []struct {
		desc      string
		oddsValue string
		betAmount int64

		expVal int64
		err    error
	}{
		{
			desc:      "upper",
			oddsValue: "+150",
			betAmount: defaultBetAmount,

			expVal: 150,
		},
		{
			desc:      "lower",
			oddsValue: "-150",
			betAmount: defaultBetAmount,

			expVal: 67,
		},
		{
			desc:      "same",
			oddsValue: "0",
			betAmount: defaultBetAmount,

			err: types.ErrMoneylineOddsCanNotBeZero,
		},
		{
			desc:      "incorrect format",
			oddsValue: "15.6",
			betAmount: defaultBetAmount,

			err: types.ErrInConvertingOddsToInt,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			payout, err := types.CalculatePayoutProfit(types.OddsType_ODD_TYPE_MONEYLINE, tc.oddsValue, sdk.NewInt(tc.betAmount))
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.expVal).Equal(payout), "expected: %d, actual: %d", tc.expVal, payout.Int64())
			}

		})
	}
}
