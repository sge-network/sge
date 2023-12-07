package keeper_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/reward/keeper"
	"github.com/sge-network/sge/x/reward/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

func getDefaultClaim(creator string) jwt.MapClaims {
	return jwt.MapClaims{
		"exp":                 time.Now().Add(time.Minute * 5).Unix(),
		"iat":                 time.Now().Unix(),
		"promoter":            creator,
		"start_ts":            time.Now().Unix(),
		"end_ts":              time.Now().Add(5 * time.Minute).Unix(),
		"category":            types.RewardType_REWARD_TYPE_SIGNUP,
		"claims_per_category": 1,
		"is_active":           true,
		"meta":                "sample campaign",
	}
}

func createCampaign(t *testing.T, k *keeper.Keeper, srv types.MsgServer, ctx sdk.Context,
	funder string, claims jwt.MapClaims,
) string {
	ticket, err := simapp.CreateJwtTicket(claims)
	require.Nil(t, err)

	expected := &types.MsgCreateCampaign{
		Creator:    funder,
		Uid:        uuid.NewString(),
		Ticket:     ticket,
		TotalFunds: sdkmath.NewInt(1000000),
	}
	_, err = srv.CreateCampaign(sdk.WrapSDKContext(ctx), expected)
	require.NoError(t, err)
	rst, found := k.GetCampaign(ctx,
		expected.Uid,
	)
	require.True(t, found)
	require.Equal(t, expected.Creator, rst.Creator)
	return expected.Uid
}

func TestMsgApplySignupReward(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	funder := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	_, err := tApp.SubaccountKeeper.CreateSubAccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	campClaims := getDefaultClaim(funder)
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

	for _, tc := range []struct {
		desc   string
		claims jwt.MapClaims
		err    error
	}{
		{
			desc: "invalid ticket",
			claims: jwt.MapClaims{
				"exp":    time.Now().Add(time.Minute * 5).Unix(),
				"iat":    time.Now().Unix(),
				"common": "invalid",
			},
			err: types.ErrInTicketVerification,
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "source id",
					Meta:      "signup reward for example user",
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     funder,
				CampaignUid: campUID,
				Ticket:      ticket,
			}
			_, err = srv.GrantReward(wctx, reward)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgApplySignupRewardSubAcc(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	funder := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	campClaims := getDefaultClaim(funder)
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

	_, err := tApp.SubaccountKeeper.CreateSubAccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	for _, tc := range []struct {
		desc   string
		claims jwt.MapClaims
		err    error
	}{
		{
			desc: "invalid ticket",
			claims: jwt.MapClaims{
				"exp":    time.Now().Add(time.Minute * 5).Unix(),
				"iat":    time.Now().Unix(),
				"common": "invalid",
			},
			err: types.ErrInTicketVerification,
		},
		{
			desc: "invalid receiver address",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  "invalid",
					SourceUID: "source id",
					Meta:      "signup reward for example user",
				},
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			desc: "subaccount not exists",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  sample.AccAddress(),
					SourceUID: "source id",
					Meta:      "signup reward for example user",
				},
			},
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "source id",
					Meta:      "signup reward for example user",
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     funder,
				CampaignUid: campUID,
				Ticket:      ticket,
			}
			_, err = srv.GrantReward(wctx, reward)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgApplySubAccFunds(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	funder := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	rewardAmount := int64(100)

	campClaims := getDefaultClaim(funder)
	campClaims["type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iat": time.Now().Unix(),
		"common": types.RewardPayloadCommon{
			Receiver:  receiverAddr,
			SourceUID: "source id",
			Meta:      "signup reward for example user",
		},
	}

	ticket, err := simapp.CreateJwtTicket(claims)
	require.Nil(t, err)

	reward := &types.MsgGrantReward{
		Uid:         uuid.NewString(),
		Creator:     funder,
		CampaignUid: campUID,
		Ticket:      ticket,
	}
	_, err = srv.GrantReward(wctx, reward)
	require.NoError(t, err)

	subAccAddr, found := tApp.SubaccountKeeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(receiverAddr))
	require.True(t, found)

	balance, found := tApp.SubaccountKeeper.GetAccountSummary(ctx, subAccAddr)
	require.True(t, found)

	require.Equal(t, rewardAmount, balance.DepositedAmount.Int64())
}
