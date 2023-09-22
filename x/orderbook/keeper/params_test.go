package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/orderbook/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := setupKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
