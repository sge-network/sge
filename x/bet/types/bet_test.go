package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

func TestCheckSettlementEligiblity(t *testing.T) {
	tcs := []struct {
		desc string
		bet  *types.Bet
		err  error
	}{
		{
			desc: "pending bet",
			bet: &types.Bet{
				Status: types.Bet_STATUS_PENDING,
			},
		},
		{
			desc: "canceled bet",
			bet: &types.Bet{
				Status: types.Bet_STATUS_CANCELED,
			},
			err: types.ErrBetIsCanceled,
		},
		{
			desc: "settled bet",
			bet: &types.Bet{
				Status: types.Bet_STATUS_SETTLED,
			},
			err: types.ErrBetIsSettled,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.bet.CheckSettlementEligiblity()
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestProcessBetResultAndStatus(t *testing.T) {
	tcs := []struct {
		desc   string
		bet    *types.Bet
		market markettypes.Market
		err    error
		result types.Bet_Result
	}{
		{
			desc: "not declared",
			market: markettypes.Market{
				Status: markettypes.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			bet: &types.Bet{},
			err: types.ErrResultNotDeclared,
		},
		{
			desc: "won",
			market: markettypes.Market{
				Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
				WinnerOddsUIDs: []string{"oddsUID"},
			},
			bet: &types.Bet{
				OddsUID: "oddsUID",
			},
			result: types.Bet_RESULT_WON,
		},
		{
			desc: "lost",
			market: markettypes.Market{
				Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
				WinnerOddsUIDs: []string{"oddsUID"},
			},
			bet:    &types.Bet{},
			result: types.Bet_RESULT_LOST,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.bet.SetResult(&tc.market)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.bet.Result, tc.result)
				require.Equal(t, tc.bet.Status, types.Bet_STATUS_RESULT_DECLARED)
			}
		})
	}
}
