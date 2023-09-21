package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/bet/types"
)

var ecceptedTruncatedValue = sdkmath.NewInt(1)

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
				types.OddsType_ODDS_TYPE_DECIMAL,
				tc.oddsValue,
				sdkmath.NewInt(tc.betAmount),
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdkmath.NewInt(tc.expVal).Equal(payoutProfit.TruncateInt()), "expected: %d, actual: %d", tc.expVal, payoutProfit)
			}

			calcBetAmount, err := types.CalculateBetAmount(
				types.OddsType_ODDS_TYPE_DECIMAL,
				tc.oddsValue,
				payoutProfit,
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdkmath.NewInt(tc.betAmount).Sub(calcBetAmount.Ceil().TruncateInt()).LT(ecceptedTruncatedValue), "expected: %d, actual: %d", tc.betAmount, calcBetAmount)
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

			expVal: 89064472,
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

			expVal: 10178796,
		},
		{
			desc:      "positive outcome 3",
			oddsValue: "1/17",
			betAmount: defaultBetAmount,

			expVal: 2095634,
		},
		{
			desc:      "positive outcome 4",
			oddsValue: "7/3",
			betAmount: defaultBetAmount,

			expVal: 83126841,
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

			err: types.ErrFractionalOddsIncorrectFormat,
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
			payoutProfit, err := types.CalculatePayoutProfit(
				types.OddsType_ODDS_TYPE_FRACTIONAL,
				tc.oddsValue,
				sdkmath.NewInt(tc.betAmount),
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdkmath.NewInt(tc.expVal).Equal(payoutProfit.TruncateInt()), "expected: %d, actual: %d", tc.expVal, payoutProfit)
			}

			calcBetAmount, err := types.CalculateBetAmount(
				types.OddsType_ODDS_TYPE_FRACTIONAL,
				tc.oddsValue,
				payoutProfit,
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdkmath.NewInt(tc.betAmount).Sub(calcBetAmount.Ceil().RoundInt()).Abs().LTE(ecceptedTruncatedValue), "expected: %d, actual: %d", tc.betAmount, calcBetAmount)
			}
		})
	}
}

func TestCalculateFractionalBetAmount(t *testing.T) {
	tcs := []struct {
		desc         string
		oddsValue    string
		payoutProfit sdk.Dec

		expVal sdk.Dec
		err    error
	}{
		{
			desc:         "positive outcome",
			oddsValue:    "2/7",
			payoutProfit: sdk.MustNewDecFromStr("5003"),

			expVal: sdk.MustNewDecFromStr("17510.5"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			calcBetAmount, err := types.CalculateBetAmount(
				types.OddsType_ODDS_TYPE_FRACTIONAL,
				tc.oddsValue,
				tc.payoutProfit,
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expVal, calcBetAmount, "expected: %d, actual: %d", tc.expVal, calcBetAmount)
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

			expVal: 160316050,
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

			err: types.ErrMoneylineOddsIncorrectFormat,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			payoutProfit, err := types.CalculatePayoutProfit(
				types.OddsType_ODDS_TYPE_MONEYLINE,
				tc.oddsValue,
				sdkmath.NewInt(tc.betAmount),
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdkmath.NewInt(tc.expVal).Equal(payoutProfit.TruncateInt()), "expected: %d, actual: %d", tc.expVal, payoutProfit)
			}

			calcBetAmount, err := types.CalculateBetAmount(
				types.OddsType_ODDS_TYPE_MONEYLINE,
				tc.oddsValue,
				payoutProfit,
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.True(t, sdkmath.NewInt(tc.betAmount).Sub(calcBetAmount.Ceil().TruncateInt()).Abs().LTE(ecceptedTruncatedValue), "expected: %d, actual: %d", tc.betAmount, calcBetAmount)
			}
		})
	}
}

func TestCalculateMoneylineBetAmount(t *testing.T) {
	tcs := []struct {
		desc         string
		oddsValue    string
		payoutProfit sdk.Dec

		expVal sdk.Dec
		err    error
	}{
		{
			desc:         "positive outcome",
			oddsValue:    "-450",
			payoutProfit: sdk.MustNewDecFromStr("5003"),

			expVal: sdk.MustNewDecFromStr("22513.5"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			calcBetAmount, err := types.CalculateBetAmount(
				types.OddsType_ODDS_TYPE_MONEYLINE,
				tc.oddsValue,
				tc.payoutProfit,
			)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expVal, calcBetAmount, "expected: %d, actual: %d", tc.expVal, calcBetAmount)
			}
		})
	}
}
