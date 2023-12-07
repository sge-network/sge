package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/sge-network/sge/x/reward/types"
)

func TestCheckBalance(t *testing.T) {
	poolBalace := int64(10000000)
	pool := types.NewPool(sdkmath.NewInt(poolBalace))
	tests := []struct {
		name  string
		spent sdkmath.Int
		err   error
	}{
		{
			name:  "not enough balance",
			spent: sdkmath.NewInt(poolBalace + 1),
			err:   types.ErrCampaignPoolBalance,
		},
		{
			name:  "valid",
			spent: sdkmath.NewInt(poolBalace),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pool.CheckBalance(tt.spent)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
