package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/cmd/sged/cmd"
	"github.com/sge-network/sge/testutil/network"
)

func TestAddGenesisAccountCmdPanic(t *testing.T) {
	userHomeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	defaultNodeHome := filepath.Join(userHomeDir, ".sge")
	panicFunc := func() {
		cmd.AddGenesisAccountCmd(defaultNodeHome)
	}

	require.NotPanics(t, panicFunc)
}

func TestSampleCMD(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx
	sec := hd.Secp256k1

	keyInfo, _, err := ctx.Keyring.NewMnemonic(
		"genUser1",
		keyring.English,
		sdk.FullFundraiserPath,
		"",
		sec,
	)
	require.NoError(t, err)

	userHomeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	defaultNodeHome := filepath.Join(userHomeDir, ".sge")

	accAddr, err := keyInfo.GetAddress()
	require.NoError(t, err)

	fields := []string{accAddr.String(), "10000000usge"}
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
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			var args []string

			args = append(args, fields...)
			args = append(args, tc.args...)
			_, err := clitestutil.ExecTestCLICmd(ctx, cmd.AddGenesisAccountCmd(defaultNodeHome), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
