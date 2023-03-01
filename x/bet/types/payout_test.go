package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

var ecceptedTruncatedValue = sdk.NewInt(7)

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

			expVal: 19594184,
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
			payout, err := types.CalculatePayoutProfit(types.OddsType_ODDS_TYPE_DECIMAL, tc.oddsValue, sdk.NewInt(tc.betAmount))
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.expVal).Equal(payout), "expected: %d, actual: %d", tc.expVal, payout.Int64())
			}

			calcBetAmount, err := types.CalculateBetAmount(types.OddsType_ODDS_TYPE_DECIMAL, tc.oddsValue, payout)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.betAmount).Sub(calcBetAmount).LT(ecceptedTruncatedValue), "expected: %d, actual: %d", tc.betAmount, calcBetAmount.Int64())
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

			expVal: 89064473,
		},
		{
			desc:      "positive outcome 1",
			oddsValue: "7/2",
			betAmount: defaultBetAmount,

			expVal: 124690261,
		},
		{
			desc:      "positive outcome 2",
			oddsValue: "2/7",
			betAmount: defaultBetAmount,

			expVal: 10178797,
		},
		{
			desc:      "positive outcome 3",
			oddsValue: "1/17",
			betAmount: defaultBetAmount,

			expVal: 2095635,
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
			payout, err := types.CalculatePayoutProfit(types.OddsType_ODDS_TYPE_FRACTIONAL, tc.oddsValue, sdk.NewInt(tc.betAmount))
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.expVal).Equal(payout), "expected: %d, actual: %d", tc.expVal, payout.Int64())
			}

			calcBetAmount, err := types.CalculateBetAmount(types.OddsType_ODDS_TYPE_FRACTIONAL, tc.oddsValue, payout)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.betAmount).Sub(calcBetAmount).Abs().LT(ecceptedTruncatedValue), "expected: %d, actual: %d", tc.betAmount, calcBetAmount.Int64())
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

			expVal: 53438683,
		},
		{
			desc:      "lower",
			oddsValue: "-450",
			betAmount: defaultBetAmount,

			expVal: 7916842,
		},
		{
			desc:      "lower",
			oddsValue: "-426",
			betAmount: defaultBetAmount,

			expVal: 8362861,
		},
		{
			desc:      "lower",
			oddsValue: "+450",
			betAmount: defaultBetAmount,

			expVal: 160316051,
		},
		{
			desc:      "lower",
			oddsValue: "+426",
			betAmount: defaultBetAmount,

			expVal: 151765861,
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
			payout, err := types.CalculatePayoutProfit(types.OddsType_ODDS_TYPE_MONEYLINE, tc.oddsValue, sdk.NewInt(tc.betAmount))
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.expVal).Equal(payout), "expected: %d, actual: %d", tc.expVal, payout.Int64())
			}

			calcBetAmount, err := types.CalculateBetAmount(types.OddsType_ODDS_TYPE_MONEYLINE, tc.oddsValue, payout)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdk.NewInt(tc.betAmount).Sub(calcBetAmount).Abs().LT(ecceptedTruncatedValue), "expected: %d, actual: %d", tc.betAmount, calcBetAmount.Int64())
			}
		})
	}
}
