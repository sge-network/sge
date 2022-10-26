package keeper_test

import (
	"testing"

	"github.com/golang-jwt/jwt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/keeper"
	"github.com/sge-network/sge/x/bet/types"
	sporteventkeeper "github.com/sge-network/sge/x/sportevent/keeper"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/require"
)

var (
	testSportEventUID = "5db09053-2901-4110-8fb5-c14e21f8d555"
	testOddsUID       = "6db09053-2901-4110-8fb5-c14e21f8d666"
	testEventOddsUIDs = []string{testOddsUID, "7db09053-2901-4110-8fb5-c14e21f8d777"}
	testCreator       string
	testBet           *types.MsgPlaceBet
	testAddSportEvent *sporteventtypes.MsgAddEvent
)

func setupKeeperAndApp(t testing.TB) (*simappUtil.TestApp, *keeper.KeeperTest, sdk.Context) {
	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	return tApp, &tApp.BetKeeper, ctx
}

func setupKeeper(t testing.TB) (*keeper.KeeperTest, sdk.Context) {
	_, k, ctx := setupKeeperAndApp(t)

	return k, ctx
}

func addSportEvent(t testing.TB, tApp *simappUtil.TestApp, ctx sdk.Context) {

	testCreator = simappUtil.TestParamUsers["user1"].Address.String()
	testAddSportEventClaim := jwt.MapClaims{
		"uid":      testSportEventUID,
		"startTS":  1111111111,
		"endTS":    9999999999,
		"oddsUIDs": testEventOddsUIDs,
		"exp":      9999999999,
		"iat":      7777777777,
	}
	testAddSportEventTicket, err := createJwtTicket(testAddSportEventClaim)
	require.Nil(t, err)

	testAddSportEvent = &sporteventtypes.MsgAddEvent{
		Creator: testCreator,
		Ticket:  testAddSportEventTicket,
	}
	wctx := sdk.WrapSDKContext(ctx)
	sporteventSrv := sporteventkeeper.NewMsgServerImpl(tApp.SporteventKeeper)
	resAddEvent, err := sporteventSrv.AddEvent(wctx, testAddSportEvent)
	require.Nil(t, err)
	require.NotNil(t, resAddEvent)
}

func placeTestBet(ctx sdk.Context, t testing.TB, tApp *simappUtil.TestApp, betUID string) {
	testCreator = simappUtil.TestParamUsers["user1"].Address.String()
	wctx := sdk.WrapSDKContext(ctx)
	betSrv := keeper.NewMsgServerImpl(tApp.BetKeeper)
	testPlaceBetClaim := jwt.MapClaims{
		"sport_event_uid": testSportEventUID,
		"value":           sdk.NewDec(10),
		"uid":             testOddsUID,
		"exp":             9999999999,
		"iat":             7777777777,
	}
	testPlaceBetTicket, err := createJwtTicket(testPlaceBetClaim)
	require.Nil(t, err)

	testBet = &types.MsgPlaceBet{
		Creator: testCreator,
		Bet: &types.BetPlaceFields{
			UID:    betUID,
			Amount: sdk.NewInt(500),
			Ticket: testPlaceBetTicket,
		},
	}
	resPlaceBet, err := betSrv.PlaceBet(wctx, testBet)
	require.Nil(t, err)
	require.NotNil(t, resPlaceBet)
}

func createJwtTicket(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claim)
	return token.SignedString(simappUtil.TestDVMPrivateKey)
}

func TestLogger(t *testing.T) {
	k, ctx := setupKeeper(t)
	logger := k.Logger(ctx)
	require.NotNil(t, logger)

	logger.Debug("test log")
}
