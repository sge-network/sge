package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/sportevent/keeper"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func createNSportEvent(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.SportEvent {
	items := make([]types.SportEvent, n)
	for i := range items {
		items[i].UID = cast.ToString(i)
		items[i].SrContributionForHouse = sdk.NewInt(0)

		keeper.SetSportEvent(ctx, items[i])
	}
	return items
}

func TestSportEventGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNSportEvent(k, ctx, 10)
	_, found := k.GetSportEvent(ctx,
		"NotExistUid",
	)
	require.False(t, found)

	for _, item := range items {
		rst, found := k.GetSportEvent(ctx,
			item.UID,
		)
		require.True(t, found)
		require.EqualValues(t, item, rst)
	}
}

func TestGetSportEvent(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNSportEvent(k, ctx, 10)

	rst, found := k.GetSportEvent(ctx, "NotExistUid")
	require.False(t, found)
	require.Equal(t, rst.UID, "")

	for _, item := range items {
		rst, found := k.GetSportEvent(ctx,
			item.UID,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(item),
			nullify.Fill(rst),
		)
	}
}

func TestSportEventRemove(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNSportEvent(k, ctx, 10)
	for _, item := range items {
		k.RemoveSportEvent(ctx,
			item.UID,
		)
		_, found := k.GetSportEvent(ctx,
			item.UID,
		)
		require.False(t, found)
	}
}

func TestSportEventGetAll(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNSportEvent(k, ctx, 10)

	sportEvents, err := k.GetSportEventAll(ctx)
	require.NoError(t, err)

	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(sportEvents),
	)
}

func TestResolveSportEvents(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		k, ctx := setupKeeper(t)
		resEventsIn := types.SportEventResolutionTicketPayload{
			UID: "NotExistUid",
		}
		_, err := k.ResolveSportEvent(ctx, &resEventsIn)
		require.Equal(t, types.ErrNoMatchingSportEvent, err)
	})

	t.Run("NotPending", func(t *testing.T) {
		k, ctx := setupKeeper(t)

		item := types.SportEvent{
			UID:    "uid",
			Status: types.SportEventStatus_SPORT_EVENT_STATUS_CANCELED,
		}
		k.SetSportEvent(ctx, item)

		resEventsIn := types.SportEventResolutionTicketPayload{
			UID: item.UID,
		}

		_, err := k.ResolveSportEvent(ctx, &resEventsIn)
		require.Equal(t, types.ErrCanNotBeAltered, err)
	})

	t.Run("Success", func(t *testing.T) {
		k, ctx := setupKeeper(t)

		item := types.SportEvent{
			UID:    "uid",
			Status: types.SportEventStatus_SPORT_EVENT_STATUS_ACTIVE,
		}
		k.SetSportEvent(ctx, item)

		resEventsIn := types.SportEventResolutionTicketPayload{
			UID:            item.UID,
			ResolutionTS:   123456,
			WinnerOddsUIDs: []string{"oddsUID1", "oddsUID2"},
			Status:         types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED,
		}
		_, err := k.ResolveSportEvent(ctx, &resEventsIn)
		require.Nil(t, err)
		val, found := k.GetSportEvent(ctx, item.UID)
		require.True(t, found)
		require.Equal(t, resEventsIn.ResolutionTS, val.ResolutionTS)
		require.Equal(t, resEventsIn.WinnerOddsUIDs, val.WinnerOddsUIDs)
		require.Equal(t, resEventsIn.Status, val.Status)
	})
}

func TestSportEventExist(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		k, ctx := setupKeeper(t)
		found := k.SportEventExists(ctx, "notExistSportEventUID")
		require.False(t, found)
	})

	t.Run("Found", func(t *testing.T) {
		k, ctx := setupKeeper(t)
		item := types.SportEvent{
			UID: "uid",
		}
		k.SetSportEvent(ctx, item)
		found := k.SportEventExists(ctx, item.UID)
		require.True(t, found)
	})
}
