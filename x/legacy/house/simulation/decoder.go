package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/x/legacy/house/types"
)

// NewDecodeStore returns a decoder function closure that unmarshal the KVPair's
// Value to the corresponding house type.
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.DepositKeyPrefix):
			var depositA, depositB types.Deposit
			cdc.MustUnmarshal(kvA.Value, &depositA)
			cdc.MustUnmarshal(kvB.Value, &depositB)
			return fmt.Sprintf("%v\n%v", depositA, depositB)
		case bytes.Equal(kvA.Key, types.WithdrawalKeyPrefix):
			var withdrawA, withdrawB types.Withdrawal
			cdc.MustUnmarshal(kvA.Value, &withdrawA)
			cdc.MustUnmarshal(kvB.Value, &withdrawB)
			return fmt.Sprintf("%v\n%v", withdrawA, withdrawB)
		default:
			panic(fmt.Sprintf(errTextInvalidHouseKey, kvA.Key))
		}
	}
}
