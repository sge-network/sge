package keeper_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestChangePubkeysVote(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

		createNActiveProposal(k, ctx, 1)
		creator := simappUtil.TestParamUsers["user1"]
		pubs, err := createNTestPubKeys(2)
		require.NoError(t, err)

		T1 := jwt.NewWithClaims(jwt.SigningMethodEdDSA, struct {
			Additions []string
			Deletions []string
			jwt.RegisteredClaims
		}{
			Additions: pubs,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
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
}
