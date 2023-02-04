package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/stretchr/testify/require"
)

func TestReserverQuery(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	reserver := &types.Reserver{
		SrPool: &types.SRPool{
			LockedAmount:   sdk.ZeroInt(),
			UnlockedAmount: sdk.NewIntFromUint64(150000000000000),
		},
	}

	k.SetReserver(ctx, *reserver)

	out, err := k.Reserver(wctx, &types.QueryReserverRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryReserverResponse{
		Reserver: reserver,
	}, out)
}
