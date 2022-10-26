package cli_test

import (
	"errors"
	"strings"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/x/mint/client/cli"
	"github.com/sge-network/sge/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestQueryInflation(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tc := range []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{},
		},
		{
			desc: "invalid",
			args: []string{"wrongArg"},
			err:  errors.New("unknown command \"wrongArg\" for \"inflation\""),
		},
		{
			desc: "invalid height",
			args: []string{"--height=1000"},
			err:  errors.New("rpc error: code = Unknown desc = cannot query with height in the future; please provide a valid height: invalid height"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := append([]string(nil), tc.args...)
			res, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryInflation(), args)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				inflation := strings.Replace(string(res.Bytes()[:]), "\n", "", 1)
				require.Equal(t, types.DefaultParams().Phases[0].Inflation.String(), inflation)
			}

		})
	}
}

func TestQueryPhaseStep(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tc := range []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{},
		},
		{
			desc: "invalid",
			args: []string{"wrongArg"},
			err:  errors.New("unknown command \"wrongArg\" for \"phase-step\""),
		},
		{
			desc: "invalid height",
			args: []string{"--height=1000"},
			err:  errors.New("rpc error: code = Unknown desc = cannot query with height in the future; please provide a valid height: invalid height"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := append([]string(nil), tc.args...)
			res, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryPhaseStep(), args)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				phaseStep := strings.Replace(string(res.Bytes()[:]), "\n", "", 1)
				require.Equal(t, "1", phaseStep)
			}

		})
	}
}

func TestQueryPhaseProvision(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tc := range []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{},
		},
		{
			desc: "invalid",
			args: []string{"wrongArg"},
			err:  errors.New("unknown command \"wrongArg\" for \"phase-provisions\""),
		},
		{
			desc: "invalid height",
			args: []string{"--height=1000"},
			err:  errors.New("rpc error: code = Unknown desc = cannot query with height in the future; please provide a valid height: invalid height"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := append([]string(nil), tc.args...)
			res, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryPhaseProvisions(), args)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				phaseProvision := strings.Replace(string(res.Bytes()[:]), "\n", "", 1)
				require.Equal(t, "0.000000000000000000", phaseProvision)
			}

		})
	}
}

func TestQueryEndPhaseStatus(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tc := range []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{},
		},
		{
			desc: "invalid",
			args: []string{"wrongArg"},
			err:  errors.New("unknown command \"wrongArg\" for \"endphase-status\""),
		},
		{
			desc: "invalid height",
			args: []string{"--height=1000"},
			err:  errors.New("rpc error: code = Unknown desc = cannot query with height in the future; please provide a valid height: invalid height"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := append([]string(nil), tc.args...)
			res, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryEndPhaseStatus(), args)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				phaseProvision := strings.Replace(string(res.Bytes()[:]), "\n", "", 1)
				require.Equal(t, "false", phaseProvision)
			}
		})
	}
}
