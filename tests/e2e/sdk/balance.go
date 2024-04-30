package sdk

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
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
