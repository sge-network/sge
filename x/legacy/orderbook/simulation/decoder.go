package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// NewDecodeStore returns a decoder function closure that unmarshal the KVPair's
// Value to the corresponding bet type.
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.OrderBookKeyPrefix):
			var marketA, marketB types.OrderBook
			cdc.MustUnmarshal(kvA.Value, &marketA)
			cdc.MustUnmarshal(kvB.Value, &marketB)
			return fmt.Sprintf("%v\n%v", marketA, marketB)
		case bytes.Equal(kvA.Key, types.OrderBookParticipationKeyPrefix):
			var marketStatsA, marketStatsB types.OrderBookParticipation
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		case bytes.Equal(kvA.Key, types.OrderBookOddsExposureKeyPrefix):
			var marketStatsA, marketStatsB types.OrderBookOddsExposure
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		case bytes.Equal(kvA.Key, types.ParticipationExposureKeyPrefix):
			var marketStatsA, marketStatsB types.ParticipationExposure
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		case bytes.Equal(kvA.Key, types.ParticipationExposureByIndexKeyPrefix):
			var marketStatsA, marketStatsB types.ParticipationExposure
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		case bytes.Equal(kvA.Key, types.HistoricalParticipationExposureKeyPrefix):
			var marketStatsA, marketStatsB types.ParticipationExposure
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		case bytes.Equal(kvA.Key, types.OrderBookStatsKeyPrefix):
			var marketStatsA, marketStatsB types.OrderBookStats
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		case bytes.Equal(kvA.Key, types.ParticipationBetPairKeyPrefix):
			var marketStatsA, marketStatsB types.ParticipationBetPair
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		default:
			panic(fmt.Sprintf(errTextInvalidOrderBookKey, kvA.Key))
		}
	}
}
