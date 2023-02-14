package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/dvm/keeper"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func createNActiveProposal(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PublicKeysChangeProposal {
	items := make([]types.PublicKeysChangeProposal, n)

	pubKeys, err := createNTestPubKeys(5)
	if err != nil {
		panic(err)
	}

	for i := range items {
		items[i].Id = uint64(i)
		items[i].Creator = simappUtil.TestParamUsers["user"+cast.ToString(i)].Address.String()
		items[i].Modifications = types.PubkeysChangeProposalPayload{Additions: pubKeys}
		items[i].StartTS = ctx.BlockTime().Unix()

		keeper.SetActivePubkeysChangeProposal(ctx, items[i])
	}
	return items
}

func TestGetActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)
	_, found := k.GetActivePubkeysChangeProposal(ctx, 5000000)
	require.False(t, found)

	for _, item := range items {
		rst, found := k.GetActivePubkeysChangeProposal(ctx, item.Id)
		require.True(t, found)
		require.EqualValues(t, item, rst)
	}
}

func TestRemoveActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)
	for _, item := range items {
		k.RemoveActiveProposal(ctx,
			item.Id,
		)
		_, found := k.GetActivePubkeysChangeProposal(ctx,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestGetAllActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)

	sportEvents, err := k.GetAllActivePubkeysChangeProposals(ctx)
	require.NoError(t, err)

	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(sportEvents),
	)
}
