package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	sgetypes "github.com/sge-network/sge/types"
	betmodulekeeper "github.com/sge-network/sge/x/bet/keeper"
	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
)

var (
	bettor1      = sample.NativeAccAddress()
	bettor1Funds = sdk.NewInt(10).Mul(micro)
)

func TestMsgServer(t *testing.T) {
	app, k, msgServer, ctx := setupMsgServerAndApp(t)

	// do subaccount creation
	require.NoError(
		t,
		simapp.FundAccount(
			app.BankKeeper,
			ctx,
			subAccFunder,
			sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).LockedBalanceDenom, subAccFunds)),
		),
	)

	_, err := msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), &types.MsgCreateSubAccount{
		Sender:          subAccFunder.String(),
		SubAccountOwner: subAccOwner.String(),
		LockedBalances: []*types.LockedBalance{
			{
				UnlockTime: time.Now().Add(24 * time.Hour),
				Amount:     subAccFunds,
			},
		},
	})
	require.NoError(t, err)

	// fund a bettor
	require.NoError(
		t,
		simapp.FundAccount(
			app.BankKeeper,
			ctx,
			bettor1,
			sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).LockedBalanceDenom, subAccFunds)),
		),
	)

	// add market
	market := addTestMarket(t, app, ctx, false)

	// do house deposit
	deposit := sdk.NewInt(1000).Mul(micro)
	depResp, err := msgServer.HouseDeposit(sdk.WrapSDKContext(ctx), houseDepositMsg(t, subAccOwner, market.UID, deposit))
	require.NoError(t, err)
	// check spend
	subBalance, exists := k.GetBalance(ctx, subAccAddr)
	require.True(t, exists)
	require.Equal(t, subBalance.SpentAmount, deposit)

	// place bet
	betMsgServer := betmodulekeeper.NewMsgServerImpl(*app.BetKeeper)
	_, err = betMsgServer.Wager(sdk.WrapSDKContext(ctx), testBet(t, bettor1, bettor1Funds))
	require.NoError(t, err)

	participateFee := app.HouseKeeper.GetHouseParticipationFee(ctx).Mul(deposit.ToDec()).TruncateInt()
	bettorFee := sdk.NewInt(100)

	t.Run("house wins", func(t *testing.T) {
		ctx, _ := ctx.CacheContext()
		app.MarketKeeper.Resolve(ctx, *market, &markettypes.MarketResolutionTicketPayload{
			UID:            market.UID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID2},
			Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		})
		err := app.BetKeeper.BatchMarketSettlements(ctx)
		require.NoError(t, err)
		err = app.OrderbookKeeper.BatchOrderBookSettlements(ctx)
		require.NoError(t, err)

		subBalance, exists := k.GetBalance(ctx, subAccAddr)
		require.True(t, exists)
		require.NoError(t, err)

		require.Equal(t, subBalance.SpentAmount.String(), sdk.ZeroInt().Add(participateFee).String())
		// check profits were forwarded to subacc owner
		ownerBalance := app.BankKeeper.GetAllBalances(ctx, subAccOwner)
		require.Equal(t,
			ownerBalance.AmountOf(k.GetParams(ctx).LockedBalanceDenom).String(),
			sdk.NewInt(10).Mul(micro).Sub(bettorFee).String())
	})

	t.Run("house loses", func(t *testing.T) {
		ctx, _ := ctx.CacheContext()
		app.MarketKeeper.Resolve(ctx, *market, &markettypes.MarketResolutionTicketPayload{
			UID:            market.UID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID1},
			Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		})
		err := app.BetKeeper.BatchMarketSettlements(ctx)
		require.NoError(t, err)
		err = app.OrderbookKeeper.BatchOrderBookSettlements(ctx)
		require.NoError(t, err)

		subBalance, exists := k.GetBalance(ctx, subAccAddr)
		require.True(t, exists)
		require.NoError(t, err)

		require.Equal(t, subBalance.SpentAmount.String(), sdk.ZeroInt().Add(participateFee).String())
		require.Equal(t, subBalance.LostAmount, bettor1Funds.Sub(bettorFee).ToDec().Mul(sdk.MustNewDecFromStr("3.2")).TruncateInt())
		// check profits were forwarded to subacc owner
		ownerBalance := app.BankKeeper.GetAllBalances(ctx, subAccOwner)
		require.Equal(t, ownerBalance.AmountOf(k.GetParams(ctx).LockedBalanceDenom), sdk.ZeroInt())
	})
	t.Run("house refund", func(t *testing.T) {
		ctx, _ := ctx.CacheContext()
		app.MarketKeeper.Resolve(ctx, *market, &markettypes.MarketResolutionTicketPayload{
			UID:            market.UID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID1},
			Status:         markettypes.MarketStatus_MARKET_STATUS_CANCELED,
		})
		err := app.BetKeeper.BatchMarketSettlements(ctx)
		require.NoError(t, err)
		err = app.OrderbookKeeper.BatchOrderBookSettlements(ctx)
		require.NoError(t, err)

		subBalance, exists := k.GetBalance(ctx, subAccAddr)
		require.True(t, exists)
		require.NoError(t, err)

		require.Equal(t, subBalance.SpentAmount, sdk.ZeroInt())
		require.Equal(t, subBalance.LostAmount, sdk.ZeroInt())
		// check profits were forwarded to subacc owner
		ownerBalance := app.BankKeeper.GetAllBalances(ctx, subAccOwner)
		require.Equal(t, ownerBalance.AmountOf(k.GetParams(ctx).LockedBalanceDenom), sdk.ZeroInt())
	})

	// TODO: not partecipated in bet fulfillment.

	t.Run("withdrawal", func(t *testing.T) {
		ctx, _ := ctx.CacheContext()
		_, err := msgServer.HouseWithdraw(sdk.WrapSDKContext(ctx), &types.MsgHouseWithdraw{Msg: houseWithdrawMsg(t, subAccOwner, deposit, depResp.Response.ParticipationIndex)})
		require.NoError(t, err)

		// do subaccount balance check
		subBalance, exists := k.GetBalance(ctx, subAccAddr)
		require.True(t, exists)

		require.Equal(t, subBalance.SpentAmount.String(), sdk.NewInt(131999680).String()) // NOTE: there was a match in the bet + participate fee
		require.Equal(t, subBalance.LostAmount.String(), sdk.ZeroInt().String())
	})
}

func houseWithdrawMsg(t testing.TB, owner sdk.AccAddress, amt sdk.Int, partecipationIndex uint64) *housetypes.MsgWithdraw {
	testKyc := &sgetypes.KycDataPayload{
		Approved: true,
		ID:       owner.String(),
	}
	ticketClaim := jwt.MapClaims{
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
		"iat":      time.Now().Unix(),
		"kyc_data": testKyc,
	}
	ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
	require.Nil(t, err)

	inputWithdraw := &housetypes.MsgWithdraw{
		Creator:            owner.String(),
		MarketUID:          testMarketUID,
		Amount:             amt,
		ParticipationIndex: partecipationIndex,
		Mode:               housetypes.WithdrawalMode_WITHDRAWAL_MODE_FULL,
		Ticket:             ticket,
	}
	return inputWithdraw
}

func houseDepositMsg(t *testing.T, owner sdk.AccAddress, uid string, amt sdk.Int) *types.MsgHouseDeposit {
	testKyc := &sgetypes.KycDataPayload{
		Approved: true,
		ID:       owner.String(),
	}
	ticketClaim := jwt.MapClaims{
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
		"iat":      time.Now().Unix(),
		"kyc_data": testKyc,
	}
	ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
	require.Nil(t, err)

	inputDeposit := &housetypes.MsgDeposit{
		Creator:   owner.String(),
		MarketUID: uid,
		Amount:    amt,
		Ticket:    ticket,
	}

	return &types.MsgHouseDeposit{
		Msg: inputDeposit,
	}
}
