package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/x/subaccount/types"
)

// NewDecodeStore returns a decoder function closure that unmarshal the KVPair's
// Value to the corresponding house type.
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.SubaccountIDPrefix):
			idA := sdk.BigEndianToUint64(kvA.Value)
			idB := sdk.BigEndianToUint64(kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)
		case bytes.Equal(kvA.Key, types.SubAccountOwnerPrefix):
			idA := sdk.AccAddress(kvA.Value)
			idB := sdk.AccAddress(kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)
		case bytes.Equal(kvA.Key, types.SubAccountOwnerReversePrefix):
			idA := sdk.AccAddress(kvA.Value)
			idB := sdk.AccAddress(kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)
		case bytes.Equal(kvA.Key, types.LockedBalancePrefix):
			var lockedBalanceA, lockedBalanceB types.LockedBalance
			cdc.MustUnmarshal(kvA.Value, &lockedBalanceA)
			cdc.MustUnmarshal(kvB.Value, &lockedBalanceB)
			return fmt.Sprintf("%v\n%v", lockedBalanceA, lockedBalanceB)
		case bytes.Equal(kvA.Key, types.AccountSummaryPrefix):
			var accSumA, accSumB types.AccountSummary
			cdc.MustUnmarshal(kvA.Value, &accSumA)
			cdc.MustUnmarshal(kvB.Value, &accSumB)
			return fmt.Sprintf("%v\n%v", accSumA, accSumB)
		default:
			panic(fmt.Sprintf(errTextInvalidHouseKey, kvA.Key))
		}
	}
}
