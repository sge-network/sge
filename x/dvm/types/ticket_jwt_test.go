package types_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestVerifyWithKey(t *testing.T) {
	Pub2, Pri2, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	bs, err := x509.MarshalPKIXPublicKey(Pub2)
	require.NoError(t, err)

	T1 := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
		PublicKeys  []string
		LeaderIndex uint32
		jwt.RegisteredClaims
	}{
		PublicKeys:  []string{string(utils.NewPubKeyMemory(bs))},
		LeaderIndex: 0,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	})
	singedT1, err := T1.SignedString(Pri2)
	require.NoError(t, err)
	ss := strings.Split(singedT1, ".")

	ticket := types.NewTestJwtToken(ss[0], ss[1], ss[2])
	require.Nil(t, err)
	t.Run("Success", func(t *testing.T) {
		_, err = ticket.VerifyJwtKey(string(utils.NewPubKeyMemory(bs)))
		require.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		invalidKey := "MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgGiRKsHZPpIWUVyHVePzoZLHLvFZ+TdnAI2Xg7WJjrJKEX5D3R5KV9uFU5lwmT09fj4BrKjwOf4Yv8+u/BJhfdsiDbkqln3FhNG1ZSxAa+9n6CKBeJku9OLpDt7olBpcydyCf8CYmTNq+YABpJbVX6iYZrbpsWK34C9fppe3rzFDAgMBAAE="
		_, err := ticket.VerifyJwtKey(invalidKey)
		require.Error(t, err)
	})
}

func TestNewTicket(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		Pub, Pri, err := ed25519.GenerateKey(rand.Reader)
		require.Nil(t, err)
		_ = Pub

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
		t.Log(err)
		it, err := types.NewJwtTicket(tkn)
		t.Log(err)
		require.Nil(t, err)
		require.NotNil(t, it)
	})

	t.Run("invalid", func(t *testing.T) {
		it, err := types.NewJwtTicket("invlaid.Token")
		require.Error(t, err)
		// require.Nil(t, it)
		_ = it
	})
}

func TestUnmarshal(t *testing.T) {
	Pub, Pri, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	_ = Pub

	// bs, err := x509.MarshalPKIXPublicKey(Pri.Public())
	// require.Nil(t, err)
	// Pbs := pem.EncodeToMemory(&pem.Block{
	// 	Type:  "PUBLIC KEY",
	// 	Bytes: bs,
	// })

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

	it, err := types.NewJwtTicket(tkn)
	require.Nil(t, err)
	require.NotNil(t, it)

	t.Run("valid", func(t *testing.T) {
		var clm struct {
			jwt.RegisteredClaims
			Title string
		}
		err := it.Unmarshal(&clm)
		require.Nil(t, err)
		require.NotEmpty(t, clm.Title)
	})
}

func TestVerify(t *testing.T) {
	Pub, Pri, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	bs, err := x509.MarshalPKIXPublicKey(Pub)
	require.Nil(t, err)
	Pbs := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bs,
	})

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

	it, err := types.NewJwtTicket(tkn)
	require.Nil(t, err)
	require.NotNil(t, it)

	t.Run("valid", func(t *testing.T) {
		err := it.Verify(string(Pbs))
		require.Nil(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		err := it.Verify("invalidPubKey")
		require.Error(t, err)
	})
	t.Run("no key", func(t *testing.T) {
		err := it.Verify("")
		require.Error(t, err)
	})
}
