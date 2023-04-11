package keeper_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"testing"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang-jwt/jwt/v4"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func createNTestPubKeys(n int) ([]string, error) {
	items := []string{}

	for i := 0; i < n; i++ {
		pub1, _, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}

		bs, err := x509.MarshalPKIXPublicKey(pub1)
		if err != nil {
			return nil, err
		}

		pb := string(utils.NewPubKeyMemory(bs))
		items = append(items, pb)
	}

	return items, nil
}

func TestChangePubkeysListProposal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		msgk, _, wctx := setupMsgServer(t)
		creator := simappUtil.TestParamUsers["user1"]
		pubs, err := createNTestPubKeys(types.MinPubKeysCount)
		require.NoError(t, err)

		T1 := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"public_keys":  pubs,
			"leader_index": 0,
			"exp":          jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedT1, err := T1.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		resp, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
			Creator: creator.Address.String(),
			Ticket:  singedT1,
		})
		require.NoError(t, err)
		require.Equal(t, true, resp.Success)
	})

	t.Run("few pubkeys", func(t *testing.T) {
		k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

		createNActiveProposal(k, ctx, 1)
		creator := simappUtil.TestParamUsers["user1"]
		pubs, err := createNTestPubKeys(4)
		require.NoError(t, err)

		proposalTicket := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"public_keys":  pubs,
			"leader_index": 0,
			"exp":          jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedProposalTicket, err := proposalTicket.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		resp, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
			Creator: creator.Address.String(),
			Ticket:  singedProposalTicket,
		})
		require.ErrorIs(t, sdkerrors.ErrInvalidRequest, err)
		require.Nil(t, resp)
	})

	t.Run("wrong index", func(t *testing.T) {
		k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

		createNActiveProposal(k, ctx, 1)
		creator := simappUtil.TestParamUsers["user1"]
		pubs, err := createNTestPubKeys(types.MinPubKeysCount)
		require.NoError(t, err)

		proposalTicket := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"public_keys":  pubs,
			"leader_index": types.MaxPubKeysCount,
			"exp":          jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedProposalTicket, err := proposalTicket.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		resp, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
			Creator: creator.Address.String(),
			Ticket:  singedProposalTicket,
		})
		require.ErrorIs(t, sdkerrors.ErrInvalidRequest, err)
		require.Nil(t, resp)
	})

	t.Run("fails", func(t *testing.T) {
		msgk, _, wctx := setupMsgServer(t)

		creator := simappUtil.TestParamUsers["user1"]
		t.Run("public keys change proposal", func(t *testing.T) {
			_, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
				Creator: creator.Address.String(),
				Ticket:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJBZGRzIjpbIi0tLS0tQkVHSU4gUFVCTElDIEtFWS0tLS0tXG5NRnd3RFFZSktvWklodmNOQVFFQkJRQURTd0F3U0FKQkFLcTU1aTg4bitDVlBTbWgwV0YrMVZHQ05nYlczMWVCXG41NmJ0SEd5WVM3QnZoMExLS2x2OTY3c3Y4WU1KTUo0eEJvc1pHa081V3lGejhGNFBEc2N5bkRzQ0F3RUFBUT09XG4tLS0tLUVORCBQVUJMSUMgS0VZLS0tLS0iXSwiZXhwIjoxNzU3NzAwMjEyfQ.PDO2Ha7Hj4SOQIbrTJeSUiKBvTnicm60VVZryoVFgfu1hCAgWqdFEYJ2aSwYU9b_O76f5AR1JDCA0roo4jJI0EbJricJLOIuPHol6Fp99rZYi4QRFt4J4ePegwMxpg-VqVqoh8ItP9GknYsOnTrtqYYCFDFjINPQ6BrUcOeEiQr5WkWOHnE_H6NfuPPOwno-pPf5lr94IIHPkWOY8VkXcFBREjJkMsUIgZmu7euI6Nim8beuReb-O3sYxtg8dCKv9vt9Bti4lAKVvESCwdcx3Yk7ARmbJvZdk4-_tKUqqJCuuDzeudyoT7Og3U1lqbzHkQyx31t8HrxRMBB19jaUTg",
			})
			require.Error(t, err)
		})
		t.Run("invalid ticekt", func(t *testing.T) {
			_, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
				Creator: creator.Address.String(),
				Ticket:  "Invalid.Ticket",
			})
			require.Error(t, err)
		})

		t.Run("unmarshal ticket", func(t *testing.T) {
			_, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
				Creator: creator.Address.String(),
				Ticket:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJBZGRzIjoiSGVsbG8iLCJleHAiOjE3NTc3MDAyMTJ9.cf9X_5pvR7Bd8Ze37u2pfeUPignG-Tg-JQayEbGUtKJcXY3ilmi4rMKaGkX44jJWbDsTlBH7zDw8Bmlr1DqmUqzEaCUg9m2gcs5qsF9dqRQsJtki0308GaTl5PnX_wlYYOSulAvYDH9o9wyGkLSSjamZZKPo3epRoefIxYeF3NMpYxZhB2zmsNqQy8oA4lxyx4ptmEv0p1nfpr9vX1P2gSojHXwXhBsFYhTsuKI2-90lPyy6Aa-u7c4kfYqrhrnJutMsyHycSe7BvC4UXUmF8cBc_uPGTY5UwISkOTAiBkxSx_n2aL-4rB8ChnrbpoyoR3HONJPBuHncVixH6nzYVQ",
			})
			require.Error(t, err)
		})
	})
}
