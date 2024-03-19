package simulation_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/simulation"
	"github.com/sge-network/sge/x/subaccount/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := app.MakeEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	subID := 100
	address := sample.NativeAccAddress()

	lockedBalance := types.LockedBalance{
		UnlockTS: 100,
		Amount:   sdkmath.NewInt(100),
	}

	accSummary := types.AccountSummary{
		DepositedAmount: sdkmath.NewInt(100),
		SpentAmount:     sdkmath.NewInt(0),
		WithdrawnAmount: sdkmath.ZeroInt(),
		LostAmount:      sdkmath.ZeroInt(),
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.SubaccountIDPrefix, Value: sdk.Uint64ToBigEndian(uint64(subID))},
			{Key: types.SubaccountOwnerPrefix, Value: address},
			{Key: types.SubaccountOwnerReversePrefix, Value: address},
			{Key: types.LockedBalancePrefix, Value: cdc.MustMarshal(&lockedBalance)},
			{Key: types.AccountSummaryPrefix, Value: cdc.MustMarshal(&accSummary)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"subID", fmt.Sprintf("%d\n%d", subID, subID)},
		{"address", fmt.Sprintf("%d\n%d", address, address)},
		{"reverseAddress", fmt.Sprintf("%d\n%d", address, address)},
		{"lockedBalances", fmt.Sprintf("%v\n%v", lockedBalance, lockedBalance)},
		{"accsum", fmt.Sprintf("%v\n%v", accSummary, accSummary)},
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
