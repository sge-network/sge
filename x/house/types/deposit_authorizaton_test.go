package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	authz "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/house/types"
)

func TestDepositGrantValidateBasic(t *testing.T) {
	tests := []struct {
		name       string
		spendLimit sdkmath.Int
		expiration time.Time
		err        error
	}{
		{
			name:       "invalid coins",
			spendLimit: sdkmath.Int{},
			expiration: time.Now().Add(5 * time.Minute),
			err:        sdkerrtypes.ErrInvalidCoins,
		},
		{
			name:       "valid",
			spendLimit: sdkmath.NewInt(10000),
			expiration: time.Now().Add(5 * time.Minute),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expTime := tt.expiration
			msgGrant, err := authz.NewMsgGrant(
				sdk.MustAccAddressFromBech32(sample.AccAddress()),
				sdk.MustAccAddressFromBech32(sample.AccAddress()),
				&types.DepositAuthorization{
					SpendLimit: tt.spendLimit,
				},
				&expTime)
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
