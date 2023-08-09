package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/house/types"
	"github.com/stretchr/testify/require"
)

func TestWithdrawGrantValidateBasic(t *testing.T) {
	tests := []struct {
		name          string
		withdrawLimit sdkmath.Int
		expiration    time.Time
		err           error
	}{
		{
			name:          "invalid coins",
			withdrawLimit: sdk.NewInt(10000),
			expiration:    time.Now().Add(5 * time.Minute),
			err:           sdkerrors.ErrInvalidCoins,
		},
		{
			name:          "valid",
			withdrawLimit: sdk.NewInt(100),
			expiration:    time.Now().Add(5 * time.Minute),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgGrant, err := authz.NewMsgGrant(
				sdk.MustAccAddressFromBech32(sample.AccAddress()),
				sdk.MustAccAddressFromBech32(sample.AccAddress()),
				&types.WithdrawAuthorization{
					WithdrawLimit: tt.withdrawLimit,
				},
				&tt.expiration)
			require.NoError(t, err)

			err = msgGrant.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
