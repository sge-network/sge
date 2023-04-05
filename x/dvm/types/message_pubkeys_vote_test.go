package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestMsgChangePubkeysVoteValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgVotePubkeysChangeRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgVotePubkeysChangeRequest{
				Creator:       "invalid_address",
				VoterKeyIndex: 0,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgVotePubkeysChangeRequest{
				Creator:       sample.AccAddress(),
				VoterKeyIndex: 0,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
