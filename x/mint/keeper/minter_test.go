package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
)

func TestTokenSupply(t *testing.T) {
	k, ctx := setupKeeper(t)
	tokenSupply := k.TokenSupply(ctx, params.DefaultBondDenom)
	require.Equal(t, int64(10020000000), tokenSupply.Int64())
}

func TestStakingTokenSupply(t *testing.T) {
	k, ctx := setupKeeper(t)
	tokenSupply := k.StakingTokenSupply(ctx)
	require.Equal(t, int64(10020000000), tokenSupply.Int64())
}

func TestBondedRatio(t *testing.T) {
	k, ctx := setupKeeper(t)
	bondedRatio := k.BondedRatio(ctx)
	expectedBondedRatio, _ := sdk.NewDecFromStr("0.000998003992015968")
	require.Equal(t, expectedBondedRatio, bondedRatio)
}

func TestMintCoins(t *testing.T) {
	k, ctx := setupKeeper(t)
	mintAmount := int64(100)
	err := k.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdkmath.NewInt(mintAmount))))
	require.NoError(t, err)

	totalSupply := k.TokenSupply(ctx, params.DefaultBondDenom)
	totalSupply = totalSupply.Add(sdkmath.NewInt(mintAmount))
	require.Equal(t, int64(10020000200), totalSupply.Int64())
}
