package cli_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/dvm/client/cli"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCmdMutation(t *testing.T) {
	net, _, pri := networkWithPublicKeys(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
	}
	Pub2, Pri2, err := ed25519.GenerateKey(rand.Reader)
	_ = Pri2
	require.NoError(t, err)
	bs, err := x509.MarshalPKIXPublicKey(Pub2)
	require.NoError(t, err)

	T1 := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
		Additions []string
		Deletions []string
		jwt.RegisteredClaims
	}{
		Additions: []string{string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: bs}))},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	})
	singedT1, err := T1.SignedString(pri)
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdMutation(), args)
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
