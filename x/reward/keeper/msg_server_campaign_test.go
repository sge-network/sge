package keeper_test

import (
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
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
	creator := simappUtil.TestParamUsers["user1"].Address.String()
	for i := 0; i < 5; i++ {
		ticketClaim := jwt.MapClaims{
			"exp":            time.Now().Add(time.Minute * 5).Unix(),
			"iat":            time.Now().Unix(),
			"funder_address": creator,
			"start_ts":       time.Now().Unix(),
			"end_ts":         time.Now().Add(5 * time.Minute).Unix(),
			"type":           types.RewardType_REWARD_TYPE_SIGNUP,
			"reward_defs": []types.Definition{
				{
					ReceiverType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
					Amount:       sdk.NewInt(100),
					DstAccType:   types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
				},
			},
			"pool_amount": sdk.NewInt(1000000),
		}
		ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
		require.Nil(t, err)

		expected := &types.MsgCreateCampaign{
			Creator: creator,
			Uid:     uuid.NewString(),
			Ticket:  ticket,
		}
		_, err = srv.CreateCampaign(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetCampaign(ctx,
			expected.Uid,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestCampaignMsgServerUpdate(t *testing.T) {
	expectedUID := uuid.NewString()

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateCampaign
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateCampaign{
				Uid: expectedUID,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateCampaign{
				Creator: sample.AccAddress(),
				Uid:     expectedUID,
			},
			err: types.ErrAuthorizationNotFound,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateCampaign{
				Uid: uuid.NewString(),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			creator := simappUtil.TestParamUsers["user1"].Address.String()

			ticketClaim := jwt.MapClaims{
				"exp":            time.Now().Add(time.Minute * 5).Unix(),
				"iat":            time.Now().Unix(),
				"funder_address": creator,
				"start_ts":       time.Now().Unix(),
				"end_ts":         time.Now().Add(5 * time.Minute).Unix(),
				"type":           types.RewardType_REWARD_TYPE_SIGNUP,
				"reward_defs": []types.Definition{
					{
						ReceiverType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
						Amount:       sdk.NewInt(100),
						DstAccType:   types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
					},
				},
				"pool_amount": sdk.NewInt(10000),
			}
			ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
			require.Nil(t, err)

			expected := &types.MsgCreateCampaign{
				Creator: creator,
				Uid:     expectedUID,
				Ticket:  ticket,
			}
			_, err = srv.CreateCampaign(wctx, expected)
			require.NoError(t, err)

			ticketClaimUpdate := jwt.MapClaims{
				"exp":    time.Now().Add(time.Minute * 5).Unix(),
				"iat":    time.Now().Unix(),
				"end_ts": time.Now().Add(5 * time.Minute).Unix(),
			}
			ticketUpdate, err := simappUtil.CreateJwtTicket(ticketClaimUpdate)
			require.Nil(t, err)
			tc.request.Ticket = ticketUpdate

			if tc.request.Creator == "" {
				tc.request.Creator = creator
			}

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
	expectedUID := uuid.NewString()

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteCampaign
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteCampaign{
				Uid: expectedUID,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteCampaign{
				Creator: sample.AccAddress(),
				Uid:     expectedUID,
			},
			err: types.ErrAuthorizationNotFound,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteCampaign{
				Uid: uuid.NewString(),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			creator := simappUtil.TestParamUsers["user1"].Address.String()

			ticketClaim := jwt.MapClaims{
				"exp":            time.Now().Add(time.Minute * 5).Unix(),
				"iat":            time.Now().Unix(),
				"funder_address": creator,
				"start_ts":       time.Now().Unix(),
				"end_ts":         time.Now().Add(5 * time.Minute).Unix(),
				"type":           types.RewardType_REWARD_TYPE_SIGNUP,
				"reward_defs": []types.Definition{
					{
						ReceiverType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
						Amount:       sdk.NewInt(100),
						DstAccType:   types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
					},
				},
				"pool_amount": sdk.NewInt(10000),
			}
			ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
			require.Nil(t, err)

			_, err = srv.CreateCampaign(wctx, &types.MsgCreateCampaign{
				Creator: creator,
				Uid:     expectedUID,
				Ticket:  ticket,
			})
			require.NoError(t, err)

			if tc.request.Creator == "" {
				tc.request.Creator = creator
			}

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
