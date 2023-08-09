package mint_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/mint/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	acc := tApp.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}
