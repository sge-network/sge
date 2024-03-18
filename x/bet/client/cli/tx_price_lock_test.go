package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/bet/client/cli"
)

func TestTXTopUpPriceLockPoolCLI(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10))).String(),
		),
	}

	t.Run("TopUp", func(t *testing.T) {
		for _, tc := range []struct {
			desc   string
			funder string
			amount string

			err  error
			code uint32
		}{
			{
				funder: sample.AccAddress(),
				amount: "1000",

				desc: "valid",
			},
			{
				funder: "invalid_address",
				amount: "1000",

				desc: "invalid address",
				err:  fmt.Errorf("any error"),
			},
			{
				funder: sample.AccAddress(),
				amount: "0",

				desc: "invalid amount",
				err:  fmt.Errorf("any error"),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					tc.funder,
					tc.amount,
				}
				args = append(args, commonArgs...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPriceLockTopUp(), args)
				if tc.err != nil {
					require.NotNil(t, err)
				} else {
					require.NoError(t, err)
					var resp sdk.TxResponse
					require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				}
			})
		}
	})
}
