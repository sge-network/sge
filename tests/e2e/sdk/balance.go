package sdk

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdknetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/sge-network/sge/app/params"
)

func GetSGEBalance(clientCtx client.Context, address string) sdkmath.Int {
	args := []string{
		address,
		fmt.Sprintf("--%s=json", flags.FlagOutput),
		fmt.Sprintf("--%s=usge", bankcli.FlagDenom),
	}
	bz, err := clitestutil.ExecTestCLICmd(clientCtx, bankcli.GetBalancesCmd(), args)
	if err != nil {
		panic(err)
	}
	respType := sdk.Coin{}
	err = clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType)
	if err != nil {
		panic(err)
	}

	return respType.Amount
}

func SendToken(val *sdknetwork.Validator, acc sdk.AccAddress, amount int) {
	// Send some funds to the new account.
	out, err := clitestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		acc,
		sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdkmath.NewInt(int64(amount)))),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, fees)).String()),
	)
	if err != nil {
		panic(err)
	}

	if !strings.Contains(out.String(), `"code":0`) {
		panic(fmt.Sprintf("failed to send tokens: %s", out.String()))
	}
}
