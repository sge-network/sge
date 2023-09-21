package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/x/bet/types"
)

// NewDecodeStore returns a decoder function closure that unmarshal the KVPair's
// Value to the corresponding bet type.
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.BetListPrefix):
			var betA, betB types.Bet
			cdc.MustUnmarshal(kvA.Value, &betA)
			cdc.MustUnmarshal(kvB.Value, &betB)
			return fmt.Sprintf("%v\n%v", betA, betB)
		case bytes.Equal(kvA.Key, types.BetIDListPrefix):
			var betIDA, betIDB types.UID2ID
			cdc.MustUnmarshal(kvA.Value, &betIDA)
			cdc.MustUnmarshal(kvB.Value, &betIDB)
			return fmt.Sprintf("%v\n%v", betIDA, betIDB)
		case bytes.Equal(kvA.Key, types.BetStatsKey):
			var betStatsA, betStatsB types.BetStats
			cdc.MustUnmarshal(kvA.Value, &betStatsA)
			cdc.MustUnmarshal(kvB.Value, &betStatsB)
			return fmt.Sprintf("%v\n%v", betStatsA, betStatsB)
		case bytes.Equal(kvA.Key, types.PendingBetListPrefix):
			var pendingBetA, pendingBetB types.PendingBet
			cdc.MustUnmarshal(kvA.Value, &pendingBetA)
			cdc.MustUnmarshal(kvB.Value, &pendingBetB)
			return fmt.Sprintf("%v\n%v", pendingBetA, pendingBetB)
		case bytes.Equal(kvA.Key, types.SettledBetListPrefix):
			var settledBetA, settleBetB types.SettledBet
			cdc.MustUnmarshal(kvA.Value, &settledBetA)
			cdc.MustUnmarshal(kvB.Value, &settleBetB)
			return fmt.Sprintf("%v\n%v", settledBetA, settleBetB)
		default:
			panic(fmt.Sprintf(errTextInvalidBetKey, kvA.Key))
		}
	}
}
