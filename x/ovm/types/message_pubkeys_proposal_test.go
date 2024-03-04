package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/ovm/types"
)

func TestMsgChangePubkeysListProposalValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgSubmitPubkeysChangeProposalRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgSubmitPubkeysChangeProposalRequest{
				Creator: "invalid_address",
			},
			err: sdkerrtypes.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgSubmitPubkeysChangeProposalRequest{
				Creator: sample.AccAddress(),
			},
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
