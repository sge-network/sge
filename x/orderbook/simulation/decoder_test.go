package simulation_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/orderbook/simulation"
	"github.com/sge-network/sge/x/orderbook/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := app.MakeEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	orderBookUID := uuid.NewString()
	orderBook := types.OrderBook{
		UID:                orderBookUID,
		ParticipationCount: 1,
		OddsCount:          1,
		Status:             types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE,
	}

	participation := types.OrderBookParticipation{
		Index:                      1,
		OrderBookUID:               orderBookUID,
		ParticipantAddress:         sample.AccAddressAsString(),
		Liquidity:                  sdk.NewInt(100),
		CurrentRoundLiquidity:      sdk.NewInt(50),
		ExposuresNotFilled:         1,
		TotalBetAmount:             sdk.NewInt(10),
		CurrentRoundTotalBetAmount: sdk.NewInt(10),
		MaxLoss:                    sdk.NewInt(10),
		CurrentRoundMaxLoss:        sdk.NewInt(10),
		CurrentRoundMaxLossOddsUID: uuid.NewString(),
		ActualProfit:               sdk.NewInt(20),
		IsSettled:                  false,
	}

	oddsExposures := types.OrderBookOddsExposure{
		OrderBookUID:     orderBookUID,
		OddsUID:          uuid.NewString(),
		FulfillmentQueue: []uint64{1},
	}

	participationExposures := types.ParticipationExposure{
		OrderBookUID:       orderBookUID,
		OddsUID:            uuid.NewString(),
		ParticipationIndex: 1,
		Exposure:           sdk.NewInt(10),
		BetAmount:          sdk.NewInt(10),
		IsFulfilled:        false,
		Round:              1,
	}

	orderBookStats := types.OrderBookStats{
		ResolvedUnsettled: []string{orderBookUID},
	}

	betPair := types.ParticipationBetPair{
		OrderBookUID:       orderBookUID,
		ParticipationIndex: 1,
		BetUID:             uuid.NewString(),
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.OrderBookKeyPrefix, Value: cdc.MustMarshal(&orderBook)},
			{Key: types.OrderBookParticipationKeyPrefix, Value: cdc.MustMarshal(&participation)},
			{Key: types.OrderBookOddsExposureKeyPrefix, Value: cdc.MustMarshal(&oddsExposures)},
			{
				Key:   types.ParticipationExposureKeyPrefix,
				Value: cdc.MustMarshal(&participationExposures),
			},
			{
				Key:   types.ParticipationExposureByIndexKeyPrefix,
				Value: cdc.MustMarshal(&participationExposures),
			},
			{
				Key:   types.HistoricalParticipationExposureKeyPrefix,
				Value: cdc.MustMarshal(&participationExposures),
			},
			{Key: types.OrderBookStatsKeyPrefix, Value: cdc.MustMarshal(&orderBookStats)},
			{Key: types.ParticipationBetPairKeyPrefix, Value: cdc.MustMarshal(&betPair)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"order_book", fmt.Sprintf("%v\n%v", orderBook, orderBook)},
		{"order_book_participation", fmt.Sprintf("%v\n%v", participation, participation)},
		{"order_book_odds_exposure", fmt.Sprintf("%v\n%v", oddsExposures, oddsExposures)},
		{
			"order_book_participation_exposure",
			fmt.Sprintf("%v\n%v", participationExposures, participationExposures),
		},
		{
			"order_book_participation_exposure_by_index",
			fmt.Sprintf("%v\n%v", participationExposures, participationExposures),
		},
		{
			"historical_order_book_participation_exposure",
			fmt.Sprintf("%v\n%v", participationExposures, participationExposures),
		},
		{"order_book_stats", fmt.Sprintf("%v\n%v", orderBookStats, orderBookStats)},
		{"order_book_participation_bet_pair", fmt.Sprintf("%v\n%v", betPair, betPair)},
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
