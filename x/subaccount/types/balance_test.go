package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/sge-network/sge/x/subaccount/types"
)

func TestLockedBalanceValidate(t *testing.T) {
	tests := []struct {
		name string
		lb   types.LockedBalance
		want error
	}{
		{
			name: "unlock time zero",
			lb: types.LockedBalance{
				UnlockTS: 0,
				Amount:   sdkmath.Int{},
			},
			want: fmt.Errorf("unlock time is zero 0"),
		},
		{
			name: "negative amount",
			lb: types.LockedBalance{
				UnlockTS: uint64(time.Now().Unix()),
				Amount:   sdkmath.NewInt(-1),
			},
			want: fmt.Errorf("amount is negative"),
		},
		{
			name: "nil amount",
			lb: types.LockedBalance{
				UnlockTS: uint64(time.Now().Unix()),
				Amount:   sdkmath.Int{},
			},
			want: fmt.Errorf("amount is nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lb.Validate()
			require.Equal(t, tt.want, got)
		})
	}
}
