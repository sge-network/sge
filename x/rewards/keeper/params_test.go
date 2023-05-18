package keeper_test

import (
	"testing"

	testkeeper "github.com/sge-network/sge/testutil/keeper"
	"github.com/sge-network/sge/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.RewardsKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
