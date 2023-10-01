package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/house/keeper"
	"github.com/sge-network/sge/x/house/types"
)

func setupMsgServerAndApp(
	t testing.TB,
) (*simapp.TestApp, *keeper.KeeperTest, types.MsgServer, sdk.Context, context.Context) {
	tApp, k, ctx := setupKeeperAndApp(t)
	return tApp, k, keeper.NewMsgServerImpl(*k), ctx, sdk.WrapSDKContext(ctx)
}
