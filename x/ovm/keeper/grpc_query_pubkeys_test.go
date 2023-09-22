package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/ovm/types"
)

func TestPubKeysList(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	t.Run("valid", func(t *testing.T) {
		resp, err := k.PubKeys(wctx, &types.QueryPubKeysRequest{})
		require.Nil(t, err)
		_ = resp
	})
	t.Run("error", func(t *testing.T) {
		resp, err := k.PubKeys(wctx, nil)
		require.Error(t, err)
		_ = resp
	})
}
