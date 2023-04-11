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
	"github.com/sge-network/sge/x/ovm/types"
	"github.com/stretchr/testify/require"
)

func TestVerifyTicket(t *testing.T) {
	k, msgk, _, wctx := setupMsgServerAndKeeper(t)

	creator := simappUtil.TestParamUsers["user1"]

	t.Run("valid", func(t *testing.T) {
		Token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
			jwt.RegisteredClaims
			Title string
		}{
			Title: "Test",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		tkn, err := Token.SignedString(&simappUtil.TestOVMPrivateKeys[0])
		require.Nil(t, err)

		err = k.VerifyTicket(wctx, tkn)
		require.Nil(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		err := k.VerifyTicket(wctx, "invalid.Token")
		require.Error(t, err)
	})
	t.Run("Verify Error", func(t *testing.T) {
		Token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
			jwt.RegisteredClaims
			Title string
		}{
			Title: "Test",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		_, intrustedPrivateKey, err := ed25519.GenerateKey(rand.Reader)
		require.Nil(t, err)
		tkn, err := Token.SignedString(intrustedPrivateKey)
		require.Nil(t, err)

		err = k.VerifyTicket(wctx, tkn)
		require.Error(t, err)
	})
	_, _, _ = msgk, wctx, creator
}

func TestVerifyTicketUnmarshal(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	creator := simappUtil.TestParamUsers["user1"]

	Pub, Pri, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	Pub2, Pri2, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	_, _ = Pub, Pub2

	var clm struct {
		jwt.RegisteredClaims
		Title string
	}
	bs, err := x509.MarshalPKIXPublicKey(Pri.Public())
	require.Nil(t, err)

	Pbs := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bs,
	})

	k.SetKeyVault(ctx, types.KeyVault{
		PublicKeys: []string{string(Pbs)},
	})

	t.Run("valid", func(t *testing.T) {
		Token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
			jwt.RegisteredClaims
			Title string
		}{
			Title: "Test",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		tkn, err := Token.SignedString(Pri)
		require.Nil(t, err)

		err = k.VerifyTicketUnmarshal(wctx, tkn, &clm)
		require.Nil(t, err)
		require.NotEmpty(t, clm.Title)
	})
	t.Run("invalid token", func(t *testing.T) {
		err = k.VerifyTicketUnmarshal(wctx, "invalid.Token", &clm)
		require.Error(t, err)
	})
	t.Run("Verify Error", func(t *testing.T) {
		Token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
			jwt.RegisteredClaims
			Title string
		}{
			Title: "Test",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		tkn, err := Token.SignedString(Pri2)
		require.Nil(t, err)

		err = k.VerifyTicketUnmarshal(wctx, tkn, &clm)
		require.Error(t, err)
	})

	t.Run("Unmarshal Error", func(t *testing.T) {
		Token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
			jwt.RegisteredClaims
			Title int
		}{
			Title: 1,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		tkn, err := Token.SignedString(Pri)
		require.Nil(t, err)

		err = k.VerifyTicketUnmarshal(wctx, tkn, &clm)
		require.Error(t, err)
	})
	_, _, _ = msgk, wctx, creator
}
