package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestPubKeysList(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	t.Run("valid", func(t *testing.T) {
		resp, err := k.ListPubKeys(wctx, &types.QueryListPubKeyAllRequest{})
		require.Nil(t, err)
		_ = resp
	})
	t.Run("error", func(t *testing.T) {
		resp, err := k.ListPubKeys(wctx, nil)
		require.Error(t, err)
		_ = resp
	})
}
