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

func TestMutation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		msgk, _, wctx := setupMsgServer(t)
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
	})

	t.Run("fails", func(t *testing.T) {
		msgk, _, wctx := setupMsgServer(t)

		creator := simappUtil.TestParamUsers["user1"]
		t.Run("mutateList", func(t *testing.T) {
			_, err := msgk.Mutation(wctx, &types.MsgMutation{
				Creator: creator.Address.String(),
				Txs:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJBZGRzIjpbIi0tLS0tQkVHSU4gUFVCTElDIEtFWS0tLS0tXG5NRnd3RFFZSktvWklodmNOQVFFQkJRQURTd0F3U0FKQkFLcTU1aTg4bitDVlBTbWgwV0YrMVZHQ05nYlczMWVCXG41NmJ0SEd5WVM3QnZoMExLS2x2OTY3c3Y4WU1KTUo0eEJvc1pHa081V3lGejhGNFBEc2N5bkRzQ0F3RUFBUT09XG4tLS0tLUVORCBQVUJMSUMgS0VZLS0tLS0iXSwiZXhwIjoxNzU3NzAwMjEyfQ.PDO2Ha7Hj4SOQIbrTJeSUiKBvTnicm60VVZryoVFgfu1hCAgWqdFEYJ2aSwYU9b_O76f5AR1JDCA0roo4jJI0EbJricJLOIuPHol6Fp99rZYi4QRFt4J4ePegwMxpg-VqVqoh8ItP9GknYsOnTrtqYYCFDFjINPQ6BrUcOeEiQr5WkWOHnE_H6NfuPPOwno-pPf5lr94IIHPkWOY8VkXcFBREjJkMsUIgZmu7euI6Nim8beuReb-O3sYxtg8dCKv9vt9Bti4lAKVvESCwdcx3Yk7ARmbJvZdk4-_tKUqqJCuuDzeudyoT7Og3U1lqbzHkQyx31t8HrxRMBB19jaUTg",
			})
			require.Error(t, err)
		})
		t.Run("Invalid Ticekt", func(t *testing.T) {
			_, err := msgk.Mutation(wctx, &types.MsgMutation{
				Creator: creator.Address.String(),
				Txs:     "Invalid.Ticket",
			})
			require.Error(t, err)
		})

		t.Run("Unmarshal Ticket", func(t *testing.T) {
			_, err := msgk.Mutation(wctx, &types.MsgMutation{
				Creator: creator.Address.String(),
				Txs:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJBZGRzIjoiSGVsbG8iLCJleHAiOjE3NTc3MDAyMTJ9.cf9X_5pvR7Bd8Ze37u2pfeUPignG-Tg-JQayEbGUtKJcXY3ilmi4rMKaGkX44jJWbDsTlBH7zDw8Bmlr1DqmUqzEaCUg9m2gcs5qsF9dqRQsJtki0308GaTl5PnX_wlYYOSulAvYDH9o9wyGkLSSjamZZKPo3epRoefIxYeF3NMpYxZhB2zmsNqQy8oA4lxyx4ptmEv0p1nfpr9vX1P2gSojHXwXhBsFYhTsuKI2-90lPyy6Aa-u7c4kfYqrhrnJutMsyHycSe7BvC4UXUmF8cBc_uPGTY5UwISkOTAiBkxSx_n2aL-4rB8ChnrbpoyoR3HONJPBuHncVixH6nzYVQ",
			})
			require.Error(t, err)
		})
	})
}
