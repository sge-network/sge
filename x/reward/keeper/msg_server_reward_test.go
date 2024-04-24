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
	sgetypes "github.com/sge-network/sge/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/reward/keeper"
	"github.com/sge-network/sge/x/reward/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

var promoterUID = uuid.NewString()

var defaultCategoryCap = []types.CategoryCap{
	{Category: types.RewardCategory_REWARD_CATEGORY_SIGNUP, CapPerAcc: 1},
}

func getDefaultClaim(creator string) jwt.MapClaims {
	return jwt.MapClaims{
		"exp":       time.Now().Add(time.Minute * 5).Unix(),
		"iat":       time.Now().Unix(),
		"promoter":  creator,
		"start_ts":  time.Now().Unix(),
		"end_ts":    time.Now().Add(5 * time.Minute).Unix(),
		"category":  types.RewardCategory_REWARD_CATEGORY_SIGNUP,
		"is_active": true,
		"meta":      "sample campaign",
	}
}

func createCampaign(t *testing.T, k *keeper.Keeper, srv types.MsgServer, ctx sdk.Context,
	promoter string, claims jwt.MapClaims, categoryCap []types.CategoryCap,
) string {
	k.SetPromoter(ctx, types.Promoter{
		Creator:   promoter,
		UID:       promoterUID,
		Addresses: []string{promoter},
		Conf: types.PromoterConf{
			CategoryCap: categoryCap,
		},
	})

	k.SetPromoterByAddress(ctx,
		types.PromoterByAddress{
			Address:     promoter,
			PromoterUID: promoterUID,
		})

	ticket, err := simapp.CreateJwtTicket(claims)
	require.Nil(t, err)

	expected := &types.MsgCreateCampaign{
		Creator:    promoter,
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

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	campClaims := getDefaultClaim(promoter)
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, promoter, campClaims, defaultCategoryCap)

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
					SourceUID: "",
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: true,
						ID:       receiverAddr,
					},
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     promoter,
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

func TestMsgApplySignupRewardWithCap(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	campClaims := getDefaultClaim(promoter)
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}
	campClaims["cap_count"] = 1

	campUID := createCampaign(t, k, srv, ctx, promoter, campClaims, []types.CategoryCap{})

	ticket, err := simapp.CreateJwtTicket(jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iat": time.Now().Unix(),
		"common": types.RewardPayloadCommon{
			Receiver:  receiverAddr,
			SourceUID: "",
			Meta:      "signup reward for sample user",
			KycData: &sgetypes.KycDataPayload{
				Approved: true,
				ID:       receiverAddr,
			},
		},
	})
	require.Nil(t, err)
	reward := &types.MsgGrantReward{
		Uid:         uuid.NewString(),
		Creator:     promoter,
		CampaignUid: campUID,
		Ticket:      ticket,
	}
	_, err = srv.GrantReward(wctx, reward)
	require.NoError(t, err)

	reward.Uid = uuid.NewString()
	_, err = srv.GrantReward(wctx, reward)
	require.ErrorContains(t, err, "maximum count cap of the campaign is reached")
}

func TestMsgApplySignupRefereeReward(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	referrer := simapp.TestParamUsers["user3"].Address.String()

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	campClaims := getDefaultClaim(promoter)
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_REFERRAL_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, promoter, campClaims, defaultCategoryCap)

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
			desc: "invalid referrer",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "",
					Meta:      "signup reward for sample user",
				},
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: referrer,
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: true,
						ID:       receiverAddr,
					},
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     promoter,
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

func TestMsgApplySignupReferrerReward(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	referee := simapp.TestParamUsers["user3"].Address.String()
	referrer := simapp.TestParamUsers["user4"].Address.String()

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	// referral signup campaign
	referralSignupCampClaims := getDefaultClaim(promoter)
	referralSignupCampClaims["reward_type"] = types.RewardType_REWARD_TYPE_REFERRAL_SIGNUP
	referralSignupCampClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	referralSignupCampClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}
	referralSignupCampUID := createCampaign(t, k, srv, ctx, promoter, referralSignupCampClaims, defaultCategoryCap)

	// referral campaign
	referralCampClaims := getDefaultClaim(promoter)
	referralCampClaims["category"] = types.RewardCategory_REWARD_CATEGORY_REFERRAL
	referralCampClaims["reward_type"] = types.RewardType_REWARD_TYPE_REFERRAL
	referralCampClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	referralCampClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}
	referralCampUID := createCampaign(t, k, srv, ctx, promoter, referralCampClaims, defaultCategoryCap)

	refereeClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iat": time.Now().Unix(),
		"common": types.RewardPayloadCommon{
			Receiver:  referee,
			SourceUID: referrer,
			Meta:      "signup reward for sample user",
			KycData: &sgetypes.KycDataPayload{
				Approved: true,
				ID:       referee,
			},
		},
	}
	refereeTicket, err := simapp.CreateJwtTicket(refereeClaims)
	require.Nil(t, err)
	reward := &types.MsgGrantReward{
		Uid:         uuid.NewString(),
		Creator:     promoter,
		CampaignUid: referralSignupCampUID,
		Ticket:      refereeTicket,
	}
	_, err = srv.GrantReward(wctx, reward)
	require.NoError(t, err)

	rewardGrant, err := k.GetRewardsOfReceiverByPromoterAndCategory(ctx, promoterUID, referee, types.RewardCategory_REWARD_CATEGORY_SIGNUP)
	require.NoError(t, err)
	require.Equal(t, types.RewardByCategory{
		UID:            reward.Uid,
		Addr:           referee,
		RewardCategory: types.RewardCategory_REWARD_CATEGORY_SIGNUP,
	}, rewardGrant[0])

	require.True(t, k.HasRewardOfReceiverByPromoter(ctx, promoterUID, referee, types.RewardCategory_REWARD_CATEGORY_SIGNUP))

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
			desc: "invalid referrer",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "",
					Meta:      "signup reward for sample user",
				},
				"referee": "invalid",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "",
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: true,
						ID:       receiverAddr,
					},
				},
				"referee": referee,
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     promoter,
				CampaignUid: referralCampUID,
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

func TestMsgApplySignupAffiliateReward(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	leadGen := simapp.TestParamUsers["user3"].Address.String()

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	campClaims := getDefaultClaim(promoter)
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_AFFILIATE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, promoter, campClaims, defaultCategoryCap)

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
			desc: "invalid lead generator",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "",
					Meta:      "signup reward for sample user",
				},
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: leadGen,
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: true,
						ID:       receiverAddr,
					},
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     promoter,
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

func TestMsgApplySignupAffiliateeReward(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	affiliatee := simapp.TestParamUsers["user3"].Address.String()
	affiliator := simapp.TestParamUsers["user4"].Address.String()

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	// referral signup campaign
	affiliateSignupCampClaims := getDefaultClaim(promoter)
	affiliateSignupCampClaims["reward_type"] = types.RewardType_REWARD_TYPE_AFFILIATE_SIGNUP
	affiliateSignupCampClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	affiliateSignupCampClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}
	affiliateSignupCampUID := createCampaign(t, k, srv, ctx, promoter, affiliateSignupCampClaims, defaultCategoryCap)

	// affiliate campaign
	affiliateCampClaims := getDefaultClaim(promoter)
	affiliateCampClaims["category"] = types.RewardCategory_REWARD_CATEGORY_AFFILIATE
	affiliateCampClaims["reward_type"] = types.RewardType_REWARD_TYPE_AFFILIATE
	affiliateCampClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	affiliateCampClaims["reward_amount"] = types.RewardAmount{
		MainAccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:      0,
	}
	affiliateCampUID := createCampaign(t, k, srv, ctx, promoter, affiliateCampClaims, defaultCategoryCap)

	affiliateClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iat": time.Now().Unix(),
		"common": types.RewardPayloadCommon{
			Receiver:  affiliatee,
			SourceUID: affiliator,
			Meta:      "signup reward for sample user",
			KycData: &sgetypes.KycDataPayload{
				Approved: true,
				ID:       affiliatee,
			},
		},
	}
	affiliateeTicket, err := simapp.CreateJwtTicket(affiliateClaims)
	require.Nil(t, err)
	reward := &types.MsgGrantReward{
		Uid:         uuid.NewString(),
		Creator:     promoter,
		CampaignUid: affiliateSignupCampUID,
		Ticket:      affiliateeTicket,
	}
	_, err = srv.GrantReward(wctx, reward)
	require.NoError(t, err)

	rewardGrant, err := k.GetRewardsOfReceiverByPromoterAndCategory(ctx, promoterUID, affiliatee, types.RewardCategory_REWARD_CATEGORY_SIGNUP)
	require.NoError(t, err)
	require.Equal(t, types.RewardByCategory{
		UID:            reward.Uid,
		Addr:           affiliatee,
		RewardCategory: types.RewardCategory_REWARD_CATEGORY_SIGNUP,
	}, rewardGrant[0])

	require.True(t, k.HasRewardOfReceiverByPromoter(ctx, promoterUID, affiliatee, types.RewardCategory_REWARD_CATEGORY_SIGNUP))

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
			desc: "invalid affiliator",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "",
					Meta:      "signup reward for sample user",
				},
				"affiliatee": "invalid",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  receiverAddr,
					SourceUID: "",
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: true,
						ID:       receiverAddr,
					},
				},
				"affiliatee": affiliatee,
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     promoter,
				CampaignUid: affiliateCampUID,
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

func TestMsgApplyBetBonus(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	bettor := simapp.TestParamUsers["user2"].Address.String()

	bet := bettypes.Bet{
		Creator:   bettor,
		UID:       uuid.NewString(),
		MarketUID: uuid.NewString(),
		Amount:    sdkmath.NewInt(301),
		Result:    bettypes.Bet_RESULT_LOST,
		Status:    bettypes.Bet_STATUS_SETTLED,
		Meta: bettypes.MetaData{
			IsMainMarket: true,
		},
	}
	tApp.BetKeeper.SetBet(ctx, bet, 1)

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, bettor, bettor, []subaccounttypes.LockedBalance{
		{
			Amount:   sdk.ZeroInt(),
			UnlockTS: uint64(ctx.BlockTime().Add(60 * time.Minute).Unix()),
		},
	})
	require.NoError(t, err)

	// referral signup campaign
	betBonusCampClaims := getDefaultClaim(promoter)
	betBonusCampClaims["category"] = types.RewardCategory_REWARD_CATEGORY_BET_DISCOUNT
	betBonusCampClaims["reward_type"] = types.RewardType_REWARD_TYPE_BET_DISCOUNT
	betBonusCampClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_PERCENTAGE
	betBonusCampClaims["reward_amount"] = types.RewardAmount{
		MainAccountPercentage: sdk.NewDecWithPrec(10, 2),
		UnlockPeriod:          0,
	}
	betBonusCampClaims["constraints"] = types.CampaignConstraints{
		MaxBetAmount: sdkmath.NewInt(300),
	}
	betBonusCampUID := createCampaign(t, k, srv, ctx, promoter, betBonusCampClaims, defaultCategoryCap)

	betBonusClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iat": time.Now().Unix(),
		"common": types.RewardPayloadCommon{
			Receiver:  bettor,
			SourceUID: "",
			Meta:      "bet bonus reward for sample user",
			KycData: &sgetypes.KycDataPayload{
				Approved: true,
				ID:       bettor,
			},
		},
		"bet_uid": bet.UID,
	}
	betBonusTicket, err := simapp.CreateJwtTicket(betBonusClaims)
	require.Nil(t, err)
	reward := &types.MsgGrantReward{
		Uid:         uuid.NewString(),
		Creator:     promoter,
		CampaignUid: betBonusCampUID,
		Ticket:      betBonusTicket,
	}
	_, err = srv.GrantReward(wctx, reward)
	require.NoError(t, err)

	rewardGrant, err := k.GetRewardsOfReceiverByPromoterAndCategory(ctx, promoterUID, bettor, types.RewardCategory_REWARD_CATEGORY_BET_DISCOUNT)
	require.NoError(t, err)
	require.Equal(t, types.RewardByCategory{
		UID:            reward.Uid,
		Addr:           bettor,
		RewardCategory: types.RewardCategory_REWARD_CATEGORY_BET_DISCOUNT,
	}, rewardGrant[0])

	require.True(t, k.HasRewardOfReceiverByPromoter(ctx, promoterUID, bettor, types.RewardCategory_REWARD_CATEGORY_BET_DISCOUNT))

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
			desc: "invalid bettor",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  bettor,
					SourceUID: "",
					Meta:      "bet bonus reward for sample user",
				},
				"bet_uid": "invalid",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"common": types.RewardPayloadCommon{
					Receiver:  bettor,
					SourceUID: "",
					Meta:      "bet bonus reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: true,
						ID:       bettor,
					},
				},
				"bet_uid": bet.UID,
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     promoter,
				CampaignUid: betBonusCampUID,
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

func TestMsgApplySignupRewardSubaccount(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	campClaims := getDefaultClaim(promoter)
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, promoter, campClaims, defaultCategoryCap)

	_, err := tApp.SubaccountKeeper.CreateSubaccount(ctx, receiverAddr, receiverAddr, []subaccounttypes.LockedBalance{
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
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: false,
						ID:       "",
						Ignore:   true,
					},
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
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: false,
						ID:       "",
						Ignore:   true,
					},
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
					Meta:      "signup reward for sample user",
					KycData: &sgetypes.KycDataPayload{
						Approved: true,
						ID:       receiverAddr,
					},
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			reward := &types.MsgGrantReward{
				Uid:         uuid.NewString(),
				Creator:     promoter,
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

func TestMsgApplySubaccountFunds(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)

	promoter := simapp.TestParamUsers["user1"].Address.String()
	receiverAddr := simapp.TestParamUsers["user2"].Address.String()

	rewardAmount := int64(100)

	campClaims := getDefaultClaim(promoter)
	campClaims["type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_type"] = types.RewardType_REWARD_TYPE_SIGNUP
	campClaims["reward_amount_type"] = types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED
	campClaims["reward_amount"] = types.RewardAmount{
		SubaccountAmount: sdkmath.NewInt(100),
		UnlockPeriod:     uint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
	}

	campUID := createCampaign(t, k, srv, ctx, promoter, campClaims, defaultCategoryCap)

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iat": time.Now().Unix(),
		"common": types.RewardPayloadCommon{
			Receiver:  receiverAddr,
			SourceUID: "source id",
			Meta:      "signup reward for sample user",
			KycData: &sgetypes.KycDataPayload{
				Approved: true,
				ID:       receiverAddr,
			},
		},
	}

	ticket, err := simapp.CreateJwtTicket(claims)
	require.Nil(t, err)

	reward := &types.MsgGrantReward{
		Uid:         uuid.NewString(),
		Creator:     promoter,
		CampaignUid: campUID,
		Ticket:      ticket,
	}
	_, err = srv.GrantReward(wctx, reward)
	require.NoError(t, err)

	subAccAddr, found := tApp.SubaccountKeeper.GetSubaccountByOwner(ctx, sdk.MustAccAddressFromBech32(receiverAddr))
	require.True(t, found)

	balance, found := tApp.SubaccountKeeper.GetAccountSummary(ctx, subAccAddr)
	require.True(t, found)

	require.Equal(t, rewardAmount, balance.DepositedAmount.Int64())
}
