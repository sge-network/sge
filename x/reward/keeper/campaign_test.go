package keeper_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/reward/keeper"
	"github.com/sge-network/sge/x/reward/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNCampaign(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Campaign {
	items := make([]types.Campaign, n)
	for i := range items {
		items[i].UID = uuid.NewString()
		items[i].Creator = sample.AccAddress()
		items[i].Promoter = sample.AccAddress()
		items[i].StartTS = uint64(time.Now().Unix())
		items[i].EndTS = uint64(time.Now().Add(5 * time.Minute).Unix())
		items[i].RewardCategory = types.RewardCategory_REWARD_CATEGORY_SIGNUP
		items[i].RewardType = types.RewardType_REWARD_TYPE_REFERRAL
		items[i].RewardAmountType = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
		items[i].IsActive = true
		items[i].Meta = "campaign " + items[i].UID
		items[i].Pool = types.Pool{Spent: sdk.ZeroInt(), Total: sdkmath.NewInt(100)}

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
			nullify.Fill(item),
			nullify.Fill(rst),
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
