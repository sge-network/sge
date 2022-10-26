package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	k, ctx := setupKeeper(t)

	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	k.SetParams(ctx, params)

	t.Run("valid", func(t *testing.T) {
		response, err := k.Params(wctx, &types.QueryParamsRequest{})
		require.NoError(t, err)
		require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
	})
	t.Run("error", func(t *testing.T) {
		_, err := k.Params(wctx, nil)
		require.Error(t, err)
	})
}
