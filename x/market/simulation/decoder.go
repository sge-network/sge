package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/x/market/types"
)

// NewDecodeStore returns a decoder function closure that unmarshal the KVPair's
// Value to the corresponding bet type.
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.MarketKeyPrefix):
			var marketA, marketB types.Market
			cdc.MustUnmarshal(kvA.Value, &marketA)
			cdc.MustUnmarshal(kvB.Value, &marketB)
			return fmt.Sprintf("%v\n%v", marketA, marketB)
		case bytes.Equal(kvA.Key, types.MarketStatsKey):
			var marketStatsA, marketStatsB types.MarketStats
			cdc.MustUnmarshal(kvA.Value, &marketStatsA)
			cdc.MustUnmarshal(kvB.Value, &marketStatsB)
			return fmt.Sprintf("%v\n%v", marketStatsA, marketStatsB)
		default:
			panic(fmt.Sprintf(errTextInvalidMarketKey, kvA.Key))
		}
	}
}
