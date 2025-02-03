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
	"github.com/sge-network/sge/x/house/simulation"
	"github.com/sge-network/sge/x/house/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := app.MakeEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	deposit := types.NewDeposit(
		sample.AccAddress(),
		sample.AccAddress(),
		uuid.NewString(),
		sdkmath.NewInt(100),
		sdkmath.NewInt(1000),
		1,
	)

	withdraw := types.Withdrawal{
		ID:                 1,
		Address:            deposit.Creator,
		MarketUID:          deposit.MarketUID,
		ParticipationIndex: 1,
		Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
		Amount:             sdkmath.NewInt(100),
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.DepositKeyPrefix, Value: cdc.MustMarshal(&deposit)},
			{Key: types.WithdrawalKeyPrefix, Value: cdc.MustMarshal(&withdraw)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"deposit", fmt.Sprintf("%v\n%v", deposit, deposit)},
		{"withdraw", fmt.Sprintf("%v\n%v", withdraw, withdraw)},
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
