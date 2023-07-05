package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestLockedBalanceValidate(t *testing.T) {
	tests := []struct {
		name string
		lb   LockedBalance
		want error
	}{
		{
			name: "unlock time zero",
			lb: LockedBalance{
				UnlockTime: time.Time{},
				Amount:     sdk.Int{},
			},
			want: fmt.Errorf("unlock time is zero"),
		},
		{
			name: "negative amount",
			lb: LockedBalance{
				UnlockTime: time.Now(),
				Amount:     sdk.NewInt(-1),
			},
			want: fmt.Errorf("amount is negative"),
		},
		{
			name: "nil amount",
			lb: LockedBalance{
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
