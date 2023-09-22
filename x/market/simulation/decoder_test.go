package simulation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/simulation"
	"github.com/sge-network/sge/x/market/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := app.MakeEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	uID := uuid.NewString()
	market := types.NewMarket(
		uID,
		sample.AccAddress(),
		cast.ToUint64(time.Now().UTC()),
		cast.ToUint64(time.Now().Add(1*time.Hour).UTC().String()),
		[]*types.Odds{
			{
				UID:  uuid.NewString(),
				Meta: "custom odds",
			},
		},
		"custom metadata",
		uID,
		types.MarketStatus_MARKET_STATUS_ACTIVE,
	)

	stats := types.MarketStats{
		ResolvedUnsettled: []string{market.UID},
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.MarketKeyPrefix, Value: cdc.MustMarshal(&market)},
			{Key: types.MarketStatsKey, Value: cdc.MustMarshal(&stats)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"market", fmt.Sprintf("%v\n%v", market, market)},
		{"market_stats", fmt.Sprintf("%v\n%v", stats, stats)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
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
