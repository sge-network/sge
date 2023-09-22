package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/x/ovm/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding bet type.
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.KeyVaultKey):
			var kvaultA, kvaultB types.KeyVault
			cdc.MustUnmarshal(kvA.Value, &kvaultA)
			cdc.MustUnmarshal(kvB.Value, &kvaultB)
			return fmt.Sprintf("%v\n%v", kvaultA, kvaultB)
		case bytes.Equal(kvA.Key, types.ProposalStatsKey):
			var statsA, statsB types.ProposalStats
			cdc.MustUnmarshal(kvA.Value, &statsA)
			cdc.MustUnmarshal(kvB.Value, &statsB)
			return fmt.Sprintf("%v\n%v", statsA, statsB)
		case bytes.Equal(kvA.Key, types.PubKeysChangeProposalListPrefix):
			var proposalA, proposalB types.PublicKeysChangeProposal
			cdc.MustUnmarshal(kvA.Value, &proposalA)
			cdc.MustUnmarshal(kvB.Value, &proposalB)
			return fmt.Sprintf("%v\n%v", proposalA, proposalB)
		default:
			panic(fmt.Sprintf(errTextInvalidOvmKey, kvA.Key))
		}
	}
}
