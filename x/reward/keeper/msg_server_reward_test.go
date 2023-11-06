package keeper_test

// import (
// 	"testing"
// 	"time"

// 	"github.com/golang-jwt/jwt"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"

// 	sdkmath "cosmossdk.io/math"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/sge-network/sge/testutil/sample"
// 	"github.com/sge-network/sge/testutil/simapp"
// 	bettypes "github.com/sge-network/sge/x/bet/types"
// 	"github.com/sge-network/sge/x/reward/keeper"
// 	"github.com/sge-network/sge/x/reward/types"
// 	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
// )

// func getDefaultClaim(creator string) jwt.MapClaims {
// 	return jwt.MapClaims{
// 		"exp":            time.Now().Add(time.Minute * 5).Unix(),
// 		"iat":            time.Now().Unix(),
// 		"promoter": creator,
// 		"start_ts":       time.Now().Unix(),
// 		"end_ts":         time.Now().Add(5 * time.Minute).Unix(),
// 		"pool_amount":    sdkmath.NewInt(1000000),
// 	}
// }

// func createCampaign(t *testing.T, k *keeper.Keeper, srv types.MsgServer, ctx sdk.Context,
// 	funder string, claims jwt.MapClaims,
// ) string {
// 	ticket, err := simapp.CreateJwtTicket(claims)
// 	require.Nil(t, err)

// 	expected := &types.MsgCreateCampaign{
// 		Creator: funder,
// 		Uid:     uuid.NewString(),
// 		Ticket:  ticket,
// 	}
// 	_, err = srv.CreateCampaign(sdk.WrapSDKContext(ctx), expected)
// 	require.NoError(t, err)
// 	rst, found := k.GetCampaign(ctx,
// 		expected.Uid,
// 	)
// 	require.True(t, found)
// 	require.Equal(t, expected.Creator, rst.Creator)
// 	return expected.Uid
// }

// func TestMsgApplySignupReward(t *testing.T) {
// 	tApp, k, ctx := setupKeeperAndApp(t)
// 	srv := keeper.NewMsgServerImpl(*k)
// 	ctx = ctx.WithBlockTime(time.Now())
// 	wctx := sdk.WrapSDKContext(ctx)

// 	funder := simapp.TestParamUsers["user1"].Address.String()
// 	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

// 	_, err := tApp.SubaccountKeeper.CreateSubAccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
// 		{
// 			Amount:   sdk.ZeroInt(),
// 			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
// 		},
// 	})
// 	require.NoError(t, err)

// 	campClaims := getDefaultClaim(funder)
// 	campClaims["type"] = types.RewardType_REWARD_TYPE_SIGNUP
// 	campClaims["reward_defs"] = []types.Definition{
// 		{
// 			RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 			RecAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_SUB,
// 			Amount:     sdkmath.NewInt(100),
// 			UnlockTS:   uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
// 		},
// 	}

// 	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

// 	for _, tc := range []struct {
// 		desc   string
// 		claims jwt.MapClaims
// 		err    error
// 	}{
// 		{
// 			desc: "invalid ticket",
// 			claims: jwt.MapClaims{
// 				"exp":      time.Now().Add(time.Minute * 5).Unix(),
// 				"iat":      time.Now().Unix(),
// 				"receiver": "invalid",
// 			},
// 			err: types.ErrInTicketVerification,
// 		},
// 		{
// 			desc: "invalid receiver type",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_REFEREE,
// 					Addr:    receiverAddr,
// 				},
// 			},
// 			err: types.ErrAccReceiverTypeNotFound,
// 		},
// 		{
// 			desc: "valid",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 					Addr:    receiverAddr,
// 				},
// 			},
// 		},
// 	} {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			ticket, err := simapp.CreateJwtTicket(tc.claims)
// 			require.Nil(t, err)
// 			reward := &types.MsgGrantReward{
// 				Creator:     funder,
// 				CampaignUid: campUID,
// 				Ticket:      ticket,
// 			}
// 			_, err = srv.GrantReward(wctx, reward)
// 			if tc.err != nil {
// 				require.ErrorContains(t, err, tc.err.Error())
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestMsgApplyAffiliationReward(t *testing.T) {
// 	k, ctx := setupKeeper(t)
// 	srv := keeper.NewMsgServerImpl(*k)
// 	ctx = ctx.WithBlockTime(time.Now())
// 	wctx := sdk.WrapSDKContext(ctx)

// 	funder := simapp.TestParamUsers["user1"].Address.String()
// 	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

// 	campClaims := getDefaultClaim(funder)
// 	campClaims["type"] = types.RewardType_REWARD_TYPE_AFFILIATION
// 	campClaims["reward_defs"] = []types.Definition{
// 		{
// 			RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 			RecAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
// 			Amount:     sdkmath.NewInt(100),
// 			UnlockTS:   0,
// 		},
// 	}

// 	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

// 	for _, tc := range []struct {
// 		desc   string
// 		claims jwt.MapClaims
// 		err    error
// 	}{
// 		{
// 			desc: "invalid ticket",
// 			claims: jwt.MapClaims{
// 				"exp":      time.Now().Add(time.Minute * 5).Unix(),
// 				"iat":      time.Now().Unix(),
// 				"receiver": "invalid",
// 			},
// 			err: types.ErrInTicketVerification,
// 		},
// 		{
// 			desc: "invalid receiver type",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_REFEREE,
// 					Addr:    receiverAddr,
// 				},
// 			},
// 			err: types.ErrAccReceiverTypeNotFound,
// 		},
// 		{
// 			desc: "valid",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 					Addr:    receiverAddr,
// 				},
// 			},
// 		},
// 	} {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			ticket, err := simapp.CreateJwtTicket(tc.claims)
// 			require.Nil(t, err)
// 			reward := &types.MsgGrantReward{
// 				Creator:     funder,
// 				CampaignUid: campUID,
// 				Ticket:      ticket,
// 			}
// 			_, err = srv.GrantReward(wctx, reward)
// 			if tc.err != nil {
// 				require.ErrorContains(t, err, tc.err.Error())
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestMsgApplyReferralReward(t *testing.T) {
// 	tApp, k, ctx := setupKeeperAndApp(t)
// 	srv := keeper.NewMsgServerImpl(*k)
// 	ctx = ctx.WithBlockTime(time.Now())
// 	wctx := sdk.WrapSDKContext(ctx)

// 	funder := simapp.TestParamUsers["user1"].Address.String()
// 	referrerAddr := simapp.TestParamUsers["user2"].Address.String()
// 	referreeAddr := simapp.TestParamUsers["user3"].Address.String()

// 	_, err := tApp.SubaccountKeeper.CreateSubAccount(ctx, referrerAddr, referrerAddr, []subaccounttypes.LockedBalance{
// 		{
// 			Amount:   sdk.ZeroInt(),
// 			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
// 		},
// 	})
// 	require.NoError(t, err)

// 	_, err = tApp.SubaccountKeeper.CreateSubAccount(ctx, referreeAddr, referreeAddr, []subaccounttypes.LockedBalance{
// 		{
// 			Amount:   sdk.ZeroInt(),
// 			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
// 		},
// 	})
// 	require.NoError(t, err)

// 	campClaims := getDefaultClaim(funder)
// 	campClaims["type"] = types.RewardType_REWARD_TYPE_REFERRAL
// 	campClaims["reward_defs"] = []types.Definition{
// 		{
// 			RecType:    types.ReceiverType_RECEIVER_TYPE_REFEREE,
// 			RecAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_SUB,
// 			Amount:     sdkmath.NewInt(100),
// 			UnlockTS:   uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
// 		},
// 		{
// 			RecType:    types.ReceiverType_RECEIVER_TYPE_REFERRER,
// 			RecAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_SUB,
// 			Amount:     sdkmath.NewInt(150),
// 			UnlockTS:   uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
// 		},
// 	}

// 	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

// 	for _, tc := range []struct {
// 		desc   string
// 		claims jwt.MapClaims
// 		err    error
// 	}{
// 		{
// 			desc: "invalid ticket",
// 			claims: jwt.MapClaims{
// 				"exp":       time.Now().Add(time.Minute * 5).Unix(),
// 				"iat":       time.Now().Unix(),
// 				"receivers": "invalid",
// 			},
// 			err: types.ErrInTicketVerification,
// 		},
// 		{
// 			desc: "invalid receiver type",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receivers": []types.Receiver{
// 					{
// 						RecType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 						Addr:    referrerAddr,
// 					},
// 					{
// 						RecType: types.ReceiverType_RECEIVER_TYPE_REFEREE,
// 						Addr:    referreeAddr,
// 					},
// 				},
// 			},
// 			err: types.ErrAccReceiverTypeNotFound,
// 		},
// 		{
// 			desc: "valid",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receivers": []types.Receiver{
// 					{
// 						RecType: types.ReceiverType_RECEIVER_TYPE_REFERRER,
// 						Addr:    referrerAddr,
// 					},
// 					{
// 						RecType: types.ReceiverType_RECEIVER_TYPE_REFEREE,
// 						Addr:    referreeAddr,
// 					},
// 				},
// 			},
// 		},
// 	} {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			ticket, err := simapp.CreateJwtTicket(tc.claims)
// 			require.Nil(t, err)
// 			reward := &types.MsgGrantReward{
// 				Creator:     funder,
// 				CampaignUid: campUID,
// 				Ticket:      ticket,
// 			}
// 			_, err = srv.GrantReward(wctx, reward)
// 			if tc.err != nil {
// 				require.ErrorContains(t, err, tc.err.Error())
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestMsgApplyNoLossBetsReward(t *testing.T) {
// 	tApp, k, ctx := setupKeeperAndApp(t)
// 	srv := keeper.NewMsgServerImpl(*k)
// 	ctx = ctx.WithBlockTime(time.Now())
// 	wctx := sdk.WrapSDKContext(ctx)

// 	funder := simapp.TestParamUsers["user1"].Address.String()
// 	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

// 	campClaims := getDefaultClaim(funder)
// 	campClaims["type"] = types.RewardType_REWARD_TYPE_NOLOSS_BETS
// 	campClaims["reward_defs"] = []types.Definition{
// 		{
// 			RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 			RecAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
// 			Amount:     sdkmath.NewInt(100),
// 			UnlockTS:   0,
// 		},
// 	}

// 	betUIDs := []string{}
// 	for i := 1; i <= 10; i++ {
// 		betUID := uuid.NewString()
// 		betUIDs = append(betUIDs, betUID)

// 		tApp.BetKeeper.SetBet(ctx, bettypes.Bet{
// 			Creator: receiverAddr,
// 			UID:     betUID,
// 			Result:  bettypes.Bet_RESULT_LOST,
// 		}, uint64(i))
// 	}

// 	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

// 	for _, tc := range []struct {
// 		desc   string
// 		claims jwt.MapClaims
// 		err    error
// 	}{
// 		{
// 			desc: "invalid ticket",
// 			claims: jwt.MapClaims{
// 				"exp":      time.Now().Add(time.Minute * 5).Unix(),
// 				"iat":      time.Now().Unix(),
// 				"receiver": "invalid",
// 			},
// 			err: types.ErrInTicketVerification,
// 		},
// 		{
// 			desc: "invalid receiver type",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_REFEREE,
// 					Addr:    receiverAddr,
// 				},
// 				"bet_uids": betUIDs,
// 			},
// 			err: types.ErrAccReceiverTypeNotFound,
// 		},
// 		{
// 			desc: "valid",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 					Addr:    receiverAddr,
// 				},
// 				"bet_uids": betUIDs,
// 			},
// 		},
// 	} {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			ticket, err := simapp.CreateJwtTicket(tc.claims)
// 			require.Nil(t, err)
// 			reward := &types.MsgGrantReward{
// 				Creator:     funder,
// 				CampaignUid: campUID,
// 				Ticket:      ticket,
// 			}
// 			_, err = srv.GrantReward(wctx, reward)
// 			if tc.err != nil {
// 				require.ErrorContains(t, err, tc.err.Error())
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestMsgApplySignupRewardSubAcc(t *testing.T) {
// 	tApp, k, ctx := setupKeeperAndApp(t)
// 	srv := keeper.NewMsgServerImpl(*k)
// 	ctx = ctx.WithBlockTime(time.Now())
// 	wctx := sdk.WrapSDKContext(ctx)

// 	funder := simapp.TestParamUsers["user1"].Address.String()
// 	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

// 	campClaims := getDefaultClaim(funder)
// 	campClaims["type"] = types.RewardType_REWARD_TYPE_SIGNUP
// 	campClaims["reward_defs"] = []types.Definition{
// 		{
// 			RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 			RecAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_SUB,
// 			Amount:     sdkmath.NewInt(100),
// 			UnlockTS:   uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
// 		},
// 	}

// 	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

// 	_, err := tApp.SubaccountKeeper.CreateSubAccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
// 		{
// 			Amount:   sdk.ZeroInt(),
// 			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
// 		},
// 	})
// 	require.NoError(t, err)

// 	for _, tc := range []struct {
// 		desc   string
// 		claims jwt.MapClaims
// 		err    error
// 	}{
// 		{
// 			desc: "invalid ticket",
// 			claims: jwt.MapClaims{
// 				"exp":      time.Now().Add(time.Minute * 5).Unix(),
// 				"iat":      time.Now().Unix(),
// 				"receiver": "invalid",
// 			},
// 			err: types.ErrInTicketVerification,
// 		},
// 		{
// 			desc: "invalid receiver type",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_REFEREE,
// 					Addr:    receiverAddr,
// 				},
// 			},
// 			err: types.ErrAccReceiverTypeNotFound,
// 		},
// 		{
// 			desc: "subaccount not exists",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 					Addr:    sample.AccAddress(),
// 				},
// 			},
// 			err: types.ErrSubAccRewardTopUp,
// 		},
// 		{
// 			desc: "valid",
// 			claims: jwt.MapClaims{
// 				"exp": time.Now().Add(time.Minute * 5).Unix(),
// 				"iat": time.Now().Unix(),
// 				"receiver": types.Receiver{
// 					RecType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 					Addr:    receiverAddr,
// 				},
// 			},
// 		},
// 	} {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			ticket, err := simapp.CreateJwtTicket(tc.claims)
// 			require.Nil(t, err)
// 			reward := &types.MsgGrantReward{
// 				Creator:     funder,
// 				CampaignUid: campUID,
// 				Ticket:      ticket,
// 			}
// 			_, err = srv.GrantReward(wctx, reward)
// 			if tc.err != nil {
// 				require.ErrorContains(t, err, tc.err.Error())
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestMsgApplySubAccFunds(t *testing.T) {
// 	tApp, k, ctx := setupKeeperAndApp(t)
// 	srv := keeper.NewMsgServerImpl(*k)
// 	ctx = ctx.WithBlockTime(time.Now())
// 	wctx := sdk.WrapSDKContext(ctx)

// 	funder := simapp.TestParamUsers["user1"].Address.String()
// 	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

// 	rewardAmount := int64(100)

// 	campClaims := getDefaultClaim(funder)
// 	campClaims["type"] = types.RewardType_REWARD_TYPE_SIGNUP
// 	campClaims["reward_defs"] = []types.Definition{
// 		{
// 			RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 			RecAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_SUB,
// 			Amount:     sdkmath.NewInt(rewardAmount),
// 			UnlockTS:   uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
// 		},
// 	}

// 	campUID := createCampaign(t, k, srv, ctx, funder, campClaims)

// 	_, err := tApp.SubaccountKeeper.CreateSubAccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
// 		{
// 			Amount:   sdk.ZeroInt(),
// 			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
// 		},
// 	})
// 	require.NoError(t, err)

// 	claims := jwt.MapClaims{
// 		"exp": time.Now().Add(time.Minute * 5).Unix(),
// 		"iat": time.Now().Unix(),
// 		"receiver": types.Receiver{
// 			RecType: types.ReceiverType_RECEIVER_TYPE_SINGLE,
// 			Addr:    receiverAddr,
// 		},
// 	}

// 	ticket, err := simapp.CreateJwtTicket(claims)
// 	require.Nil(t, err)

// 	reward := &types.MsgGrantReward{
// 		Creator:     funder,
// 		CampaignUid: campUID,
// 		Ticket:      ticket,
// 	}
// 	_, err = srv.GrantReward(wctx, reward)
// 	require.NoError(t, err)

// 	subAccAddr, found := tApp.SubaccountKeeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(receiverAddr))
// 	require.True(t, found)

// 	balance, found := tApp.SubaccountKeeper.GetBalance(ctx, subAccAddr)
// 	require.True(t, found)

// 	require.Equal(t, rewardAmount, balance.DepositedAmount.Int64())
// }
