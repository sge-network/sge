package keeper_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestQueryPublicKeys(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	creator := simappUtil.TestParamUsers["user1"]
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
	singedT1, err := T1.SignedString(simappUtil.TestDVMPrivateKey)
	require.NoError(t, err)

	resp, err := msgk.Mutation(wctx, &types.MsgMutation{
		Creator: creator.Address.String(),
		Txs:     singedT1,
	})
	require.Nil(t, err)
	require.Equal(t, true, resp.Success)

	l, found := k.GetPublicKeys(ctx)
	require.True(t, found)
	require.Greater(t, len(l.List), 0)
}

func TestSetKeys(t *testing.T) {
	k, ctx := setupKeeper(t)
	t.Run("valid", func(t *testing.T) {
		(*k).SetPublicKeys(ctx, types.PublicKeys{
			List: []string{"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+9wlxVu9a8lzUO2kcFLu\nUBIuV0+DpUdgEmsyQXr4y65sPSx/XjbK3GSZS1fB4irYPPG8EPHa6Z9KwWJLrTBr\nHayQcUBV5GQPf7nDktCkljYEBRmJZ+x3tlTf2kyKf3JMPAYgSFcs792dMpx8EiuE\n683QzUyeCutmiSWj1e7/IR9tjD4X/XFGkLES6wtqpQpOsL10z3hZllQEqZif8pDZ\nZcDvF97dg0l+JIWW3jBINL/UzuBRmdtDMuS1d57bpaMNb7L9HLUDBiwlZTGhs1+v\n9eTMY6IEdIzQ6M1KTFDeLYdnpGWP0ttBpt7SesLNpsKStbZ7QkbNtzlkTN8eJ6qu\nJQIDAQAB\n-----END PUBLIC KEY-----"},
		})
	})
}
