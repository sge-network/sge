package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

func TestBetConstraintValidation(t *testing.T) {
	param := types.DefaultParams()

	tests := []struct {
		name           string
		betConstraints *types.MarketBetConstraints
		err            error
	}{
		{
			name: "valid",
			betConstraints: &types.MarketBetConstraints{
				MinAmount: param.MinBetAmount,
				BetFee:    param.MaxBetFee,
			},
		},
		{
			name: "min bet amount",
			betConstraints: &types.MarketBetConstraints{
				MinAmount: param.MinBetAmount.Sub(sdk.OneInt()),
				BetFee:    sdk.NewInt(10),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "negative min bet amount",
			betConstraints: &types.MarketBetConstraints{
				MinAmount: sdk.NewInt(-1),
				BetFee:    sdk.NewInt(10),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "min bet fee exceeded",
			betConstraints: &types.MarketBetConstraints{
				MinAmount: param.MinBetAmount,
				BetFee:    param.MinBetFee.Sub(sdk.OneInt()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "max bet fee exceeded",
			betConstraints: &types.MarketBetConstraints{
				MinAmount: param.MinBetAmount,
				BetFee:    param.MaxBetFee.Add(sdk.OneInt()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.betConstraints.Validate(&param)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
