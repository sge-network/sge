package types_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestMsgChangePubkeysVoteValidateBasic(t *testing.T) {
	pubkey, _, _ := ed25519.GenerateKey(rand.Reader)
	bs, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		panic(err)
	}
	pubKeyStr := string(utils.NewPubKeyMemory(bs))

	tests := []struct {
		name string
		msg  types.MsgVotePubkeysChangeRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgVotePubkeysChangeRequest{
				Creator:   "invalid_address",
				PublicKey: pubKeyStr,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgVotePubkeysChangeRequest{
				Creator:   sample.AccAddress(),
				PublicKey: pubKeyStr,
			},
		},
		{
			name: "empty public key",
			msg: types.MsgVotePubkeysChangeRequest{
				Creator:   sample.AccAddress(),
				PublicKey: "",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
