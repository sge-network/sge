package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/reward/keeper"
	"github.com/sge-network/sge/x/reward/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNCampaign(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Campaign {
	items := make([]types.Campaign, n)
	for i := range items {
		items[i].UID = uuid.NewString()

		keeper.SetCampaign(ctx, items[i])
	}
	return items
}

func TestCampaignGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNCampaign(k, ctx, 10)
	for _, item := range items {
		rst, found := k.GetCampaign(ctx,
			item.UID,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestCampaignRemove(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNCampaign(k, ctx, 10)
	for _, item := range items {
		k.RemoveCampaign(ctx,
			item.UID,
		)
		_, found := k.GetCampaign(ctx,
			item.UID,
		)
		require.False(t, found)
	}
}

func TestCampaignGetAll(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNCampaign(k, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(k.GetAllCampaign(ctx)),
	)
}
