package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/ovm/types"
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
