package sdk

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdknetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz/client/cli"
	authzclitestutil "github.com/cosmos/cosmos-sdk/x/authz/client/testutil"
	"github.com/sge-network/sge/app/params"
)

func SetGenericAuthorization(val *sdknetwork.Validator, grantee sdk.AccAddress, typeMsg string) {
	twoHours := time.Now().Add(time.Minute * 120).Unix()

	// Send some funds to the new account.
	out, err := authzclitestutil.CreateGrant(
		val.ClientCtx,
		[]string{
			grantee.String(),
			"generic",
			fmt.Sprintf("--%s=%s", cli.FlagMsgType, typeMsg),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			fmt.Sprintf("--%s=%d", cli.FlagExpiration, twoHours),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, fees)).String()),
		},
	)
	if err != nil {
		panic(err)
	}

	if !strings.Contains(out.String(), `"code":0`) {
		panic(fmt.Sprintf("failed to send tokens: %s", out.String()))
	}
}
