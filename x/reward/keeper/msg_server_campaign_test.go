package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/reward/keeper"
	"github.com/sge-network/sge/x/reward/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCampaignMsgServerCreate(t *testing.T) {
	k, ctx := setupKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateCampaign{Creator: creator,
			Uid: uuid.NewString(),
		}
		_, err := srv.CreateCampaign(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetCampaign(ctx,
			expected.Uid,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestCampaignMsgServerUpdate(t *testing.T) {
	creator := "A"
	expectedUID := uuid.NewString()

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateCampaign
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateCampaign{Creator: creator,
				Uid: expectedUID,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateCampaign{Creator: "B",
				Uid: expectedUID,
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateCampaign{Creator: creator,
				Uid: uuid.NewString(),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateCampaign{Creator: creator,
				Uid: expectedUID,
			}
			_, err := srv.CreateCampaign(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateCampaign(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetCampaign(ctx,
					expected.Uid,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestCampaignMsgServerDelete(t *testing.T) {
	creator := "A"
	expectedUID := uuid.NewString()

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteCampaign
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteCampaign{Creator: creator,
				Uid: expectedUID,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteCampaign{Creator: "B",
				Uid: expectedUID,
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteCampaign{Creator: creator,
				Uid: uuid.NewString(),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateCampaign(wctx, &types.MsgCreateCampaign{Creator: creator,
				Uid: expectedUID,
			})
			require.NoError(t, err)
			_, err = srv.DeleteCampaign(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetCampaign(ctx,
					tc.request.Uid,
				)
				require.False(t, found)
			}
		})
	}
}
