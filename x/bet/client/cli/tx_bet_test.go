package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/x/bet/client/cli"
	"github.com/sge-network/sge/x/bet/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestTXBetCLI(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	t.Run("Place bet", func(t *testing.T) {
		for _, tc := range []struct {
			desc     string
			uid      string
			amount   string
			oddsType string
			ticket   string

			err  error
			code uint32
		}{
			{
				uid:      "6e31c60f-2025-48ce-ae79-1dc110f16355",
				amount:   "555",
				oddsType: "1",
				ticket:   "ticket",

				desc: "valid",
			},
			{
				uid:      "invalidUID",
				amount:   "555",
				oddsType: "2",
				ticket:   "ticket",

				desc: "validation failed",
				err:  fmt.Errorf("any error"),
			},
			{
				uid:      "6e31c60f-2025-48ce-ae79-1dc110f16355",
				amount:   "invalidAmount",
				oddsType: "1",
				ticket:   "ticket",

				desc: "invalid amuont",
				err:  fmt.Errorf("any error"),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					tc.uid,
					tc.amount,
					tc.oddsType,
					tc.ticket,
				}
				args = append(args, commonArgs...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPlaceBet(), args)
				if tc.err != nil {
					require.NotNil(t, err)
				} else {
					require.NoError(t, err)
					var resp sdk.TxResponse
					require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
					//require.Equal(t, tc.code, resp.Code)
				}
			})
		}
	})

	t.Run("Settle bet", func(t *testing.T) {
		for _, tc := range []struct {
			desc   string
			betUID string

			err    error
			errMsg string
			code   uint32
		}{
			{
				betUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",

				desc: "valid",
			},
			{
				betUID: "",

				desc:   "validation failed",
				err:    fmt.Errorf("any error"),
				errMsg: types.ErrInvalidBetUID.Error(),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					ctx.GetFromAddress().String(),
					tc.betUID,
				}
				args = append(args, commonArgs...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSettleBet(), args)
				if tc.err != nil {
					require.NotNil(t, err)
					if tc.errMsg != "" {
						require.Equal(t, tc.errMsg, err.Error())
					}
				} else {
					require.NoError(t, err)
					var resp sdk.TxResponse
					require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
					//require.Equal(t, tc.code, resp.Code)
				}
			})
		}
	})
}
