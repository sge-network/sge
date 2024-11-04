package simulation_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/bet/simulation"
	"github.com/sge-network/sge/x/bet/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := app.MakeEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	bet := types.NewBet(
		sample.AccAddress(),
		&types.WagerProps{
			UID:    uuid.NewString(),
			Amount: sdkmath.NewInt(10),
			Ticket: "",
		},
		&types.BetOdds{
			UID:               uuid.NewString(),
			MarketUID:         uuid.NewString(),
			Value:             "100",
			MaxLossMultiplier: sdkmath.LegacyNewDec(1),
		},
		types.MetaData{
			SelectedOddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
			SelectedOddsValue: "1.5",
		},
	)

	betUID := types.UID2ID{
		UID: bet.UID,
		ID:  1,
	}

	betStats := types.BetStats{
		Count: 1,
	}

	pendingBet := types.PendingBet{
		UID:     bet.UID,
		Creator: bet.Creator,
	}

	settledBet := types.SettledBet{
		UID:           bet.UID,
		BettorAddress: bet.Creator,
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.BetListPrefix, Value: cdc.MustMarshal(bet)},
			{Key: types.BetIDListPrefix, Value: cdc.MustMarshal(&betUID)},
			{Key: types.BetStatsKey, Value: cdc.MustMarshal(&betStats)},
			{Key: types.PendingBetListPrefix, Value: cdc.MustMarshal(&pendingBet)},
			{Key: types.SettledBetListPrefix, Value: cdc.MustMarshal(&settledBet)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"bet", fmt.Sprintf("%v\n%v", *bet, *bet)},
		{"bet_uid", fmt.Sprintf("%v\n%v", betUID, betUID)},
		{"bet_stats", fmt.Sprintf("%v\n%v", betStats, betStats)},
		{"pending_bets", fmt.Sprintf("%v\n%v", pendingBet, pendingBet)},
		{"settled_bets", fmt.Sprintf("%v\n%v", settledBet, settledBet)},
		{"other", ""},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
