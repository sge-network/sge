package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/subaccount/types"
)

func TestKeeper_GetParams(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	params := k.GetParams(ctx)
	require.Equal(t, params, types.DefaultParams())
}
