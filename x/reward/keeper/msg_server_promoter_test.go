package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/reward/keeper"
	"github.com/sge-network/sge/x/reward/types"
)

func TestSetPromoterConfig(t *testing.T) {
	k, ctx := setupKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	ctx = ctx.WithBlockTime(time.Now())
	wctx := sdk.WrapSDKContext(ctx)
	promoter := simapp.TestParamUsers["user1"].Address.String()
	promoterUID := uuid.NewString()
	k.SetPromoter(ctx, types.Promoter{
		Creator: promoter,
		UID:     promoterUID,
		Addresses: []string{
			promoter,
		},
		Conf: types.PromoterConf{
			CategoryCap: []types.CategoryCap{
				{
					Category:  types.RewardCategory_REWARD_CATEGORY_SIGNUP,
					CapPerAcc: 1,
				},
			},
		},
	})

	for _, tc := range []struct {
		desc   string
		claims jwt.MapClaims
		uid    string
		err    error
	}{
		{
			desc: "invalid ticket",
			claims: jwt.MapClaims{
				"exp":  time.Now().Add(time.Minute * 5).Unix(),
				"iat":  time.Now().Unix(),
				"conf": "invalid",
			},
			uid: promoterUID,
			err: types.ErrInTicketVerification,
		},
		{
			desc: "duplicate cap category",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"conf": types.PromoterConf{
					CategoryCap: []types.CategoryCap{
						{Category: types.RewardCategory_REWARD_CATEGORY_SIGNUP, CapPerAcc: 1},
						{Category: types.RewardCategory_REWARD_CATEGORY_SIGNUP, CapPerAcc: 1},
					},
				},
			},
			uid: promoterUID,
			err: types.ErrDuplicateCategoryInConf,
		},
		{
			desc: "low cap category",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"conf": types.PromoterConf{
					CategoryCap: []types.CategoryCap{
						{Category: types.RewardCategory_REWARD_CATEGORY_SIGNUP, CapPerAcc: 0},
						{Category: types.RewardCategory_REWARD_CATEGORY_AFFILIATE, CapPerAcc: 1},
					},
				},
			},
			uid: promoterUID,
			err: types.ErrCategoryCapShouldBePos,
		},
		{
			desc: "valid",
			claims: jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"conf": types.PromoterConf{
					CategoryCap: []types.CategoryCap{
						{Category: types.RewardCategory_REWARD_CATEGORY_SIGNUP, CapPerAcc: 1},
					},
				},
			},
			uid: promoterUID,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ticket, err := simapp.CreateJwtTicket(tc.claims)
			require.Nil(t, err)
			conf := &types.MsgSetPromoterConf{
				Uid:     tc.uid,
				Creator: promoter,
				Ticket:  ticket,
			}
			_, err = srv.SetPromoterConf(wctx, conf)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
