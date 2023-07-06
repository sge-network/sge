package types_test

import (
	"fmt"
	"github.com/sge-network/sge/x/subaccount/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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
				UnlockTime: time.Time{},
				Amount:     sdk.Int{},
			},
			want: fmt.Errorf("unlock time is zero"),
		},
		{
			name: "negative amount",
			lb: types.LockedBalance{
				UnlockTime: time.Now(),
				Amount:     sdk.NewInt(-1),
			},
			want: fmt.Errorf("amount is negative"),
		},
		{
			name: "nil amount",
			lb: types.LockedBalance{
				UnlockTime: time.Now(),
				Amount:     sdk.Int{},
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
