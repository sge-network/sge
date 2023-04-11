package cli_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt/v4"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/dvm/client/cli"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestCmdChangePubkeysListProposal(t *testing.T) {
	net, _ := networkWithPublicKeys(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	var pubs []string
	for i := 0; i < types.MinPubKeysCount; i++ {
		pub, _, err := ed25519.GenerateKey(rand.Reader)
		require.NoError(t, err)
		bs, err := x509.MarshalPKIXPublicKey(pub)
		require.NoError(t, err)
		pubs = append(pubs, string(utils.NewPubKeyMemory(bs)))
	}

	t1 := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"public_keys":  pubs,
		"leader_index": 0,
		"exp":          jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	})
	singedT1, err := t1.SignedString(simappUtil.TestDVMPrivateKeys[0])
	require.NoError(t, err)

	TestCases := []struct {
		desc   string
		ticket string

		err  error
		code uint32
	}{
		{
			desc:   "success",
			ticket: singedT1,
			err:    nil,
			code:   0,
		},
	}
	for _, tc := range TestCases {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.ticket,
			}
			args = append(args, commonArgs...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdChangePubkeysListProposal(), args)
			if tc.err != nil {
				require.Error(t, err, "")
			} else {
				require.NoError(t, err)

				var resp sdk.TxResponse

				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				fmt.Println(resp)
				require.True(t, resp.Code == tc.code)
			}
		})
	}
}
