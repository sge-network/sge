package keeper_test

import (
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_GetParams(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	params := k.GetParams(ctx)
	require.Equal(t, params, types.DefaultParams())
}
