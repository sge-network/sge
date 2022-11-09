package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	"github.com/stretchr/testify/require"
)

func TestTokenSupply(t *testing.T) {
	k, ctx := setupKeeper(t)
	tokenSupply := k.TokenSupply(ctx, params.DefaultBondDenom)
	require.Equal(t, int64(150010020002000), tokenSupply.Int64())
}

func TestStakingTokenSupply(t *testing.T) {
	k, ctx := setupKeeper(t)
	tokenSupply := k.StakingTokenSupply(ctx)
	require.Equal(t, int64(150010020002000), tokenSupply.Int64())
}

func TestBondedRatio(t *testing.T) {
	k, ctx := setupKeeper(t)
	bondedRatio := k.BondedRatio(ctx)
	expectedBondedRatio, _ := sdk.NewDecFromStr("0.000000066662213629")
	require.Equal(t, expectedBondedRatio, bondedRatio)
}

func TestMintCoins(t *testing.T) {
	k, ctx := setupKeeper(t)
	mintAmount := int64(100)
	k.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(mintAmount))))
	totalSupply := k.TokenSupply(ctx, params.DefaultBondDenom)
	totalSupply = totalSupply.Add(sdk.NewInt(mintAmount))
	require.Equal(t, int64(150010020002200), totalSupply.Int64())
}
