package rewards_test

import (
	"testing"

	keepertest "github.com/sge-network/sge/testutil/keeper"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/rewards"
	"github.com/sge-network/sge/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.RewardsKeeper(t)
	rewards.InitGenesis(ctx, *k, genesisState)
	got := rewards.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
